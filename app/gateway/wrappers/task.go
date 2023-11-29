package wrappers

import (
	// 引入所需的包
	"context" // 用于传递跨API边界的上下文
	"strconv" // 用于字符串和其他类型的转换

	"github.com/afex/hystrix-go/hystrix" // Hystrix断路器库
	"go-micro.dev/v4/client"             // go-micro的客户端包

	"github.com/CocaineCong/micro-todoList/idl/pb" // 项目特定的协议缓冲区包
)

// NewTask 创建新任务的函数
func NewTask(id uint64, name string) *pb.TaskModel {
	// 构造并返回一个指向TaskModel的指针
	return &pb.TaskModel{
		Id:         id,     // 任务ID
		Title:      name,   // 任务名称
		Content:    "响应超时", // 默认内容
		StartTime:  1000,   // 开始时间（示例值）
		EndTime:    1000,   // 结束时间（示例值）
		Status:     0,      // 状态（示例值）
		CreateTime: 1000,   // 创建时间（示例值）
		UpdateTime: 1000,   // 更新时间（示例值）
	}
}

// DefaultTasks 降级函数，用于熔断时提供默认响应
func DefaultTasks(resp interface{}) {
	// 创建TaskModel切片
	models := make([]*pb.TaskModel, 0)
	var i uint64
	for i = 0; i < 10; i++ {
		// 填充默认任务数据
		models = append(models, NewTask(i, "降级备忘录"+strconv.Itoa(20+int(i))))
	}
	// 将默认响应设置为任务列表响应
	result := resp.(*pb.TaskListResponse)
	result.TaskList = models
}

// TaskWrapper 结构体，包含客户端，并实现断路器逻辑
type TaskWrapper struct {
	client.Client // 嵌入go-micro客户端
}

// Call 重写Client的Call方法以实现断路器逻辑
func (wrapper *TaskWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	// 构造命令名称
	cmdName := req.Service() + "." + req.Endpoint()
	// 配置熔断器参数
	config := hystrix.CommandConfig{
		Timeout:                3000,
		RequestVolumeThreshold: 20,   // 熔断器请求阈值，默认20，意思是有20个请求才能进行错误百分比计算
		ErrorPercentThreshold:  50,   // 错误百分比，当错误超过百分比时，直接进行降级处理，直至熔断器再次 开启，默认50%
		SleepWindow:            5000, // 过多长时间，熔断器再次检测是否开启，单位毫秒ms（默认5秒）
	}
	// 配置熔断器
	hystrix.ConfigureCommand(cmdName, config)
	// 执行带有熔断器的请求
	return hystrix.Do(cmdName, func() error {
		return wrapper.Client.Call(ctx, req, rsp)
	}, func(err error) error {
		// 熔断时执行的降级函数
		DefaultTasks(rsp)
		return err
	})
}

// NewTaskWrapper 初始化Wrapper的工厂函数
func NewTaskWrapper(c client.Client) client.Client {
	return &TaskWrapper{c} // 返回TaskWrapper实例
}
