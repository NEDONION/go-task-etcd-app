package rpc

import (
	"context"
	"errors"

	"github.com/CocaineCong/micro-todoList/idl/pb"
	"github.com/CocaineCong/micro-todoList/pkg/e"
)

// UserLogin 用户登陆
func UserLogin(ctx context.Context, req *pb.UserRequest) (resp *pb.UserDetailResponse, err error) {
	resp, err = UserService.UserLogin(ctx, req)
	if err != nil {
		return
	}

	if resp.Code != e.SUCCESS {
		err = errors.New(e.GetMsg(int(resp.Code)))
		return
	}

	return
}

// UserRegister 用户注册
// *pb.UserRequest = 指向pb.UserRequest 结构体的指针
// context.Context = 标准库中的上下文，用来控制请求的超时时间、截止时间、传递请求的值等

/*
为什么要使用*？指向该类型的指针。

1. 性能优化：对于较大的数据结构，使用指针可以避免在函数调用时复制整个数据结构。这可以减少内存使用并提高程序性能。
2. 允许修改数据：当你通过指针传递变量时，函数可以修改原始数据。如果你传递一个非指针变量（即值传递），函数将获得该变量的副本，因此无法修改原始数据。
3. 一致性：在某些场景中，比如使用接口或者在结构体中嵌入其他结构体时，使用指针可以保持类型的一致性。
4. 空值处理：指针可以是nil，这对于表示“没有值”或“零值”是有用的，特别是在错误处理和可选字段中。

*/

func UserRegister(ctx context.Context, req *pb.UserRequest) (resp *pb.UserDetailResponse, err error) {
	resp, err = UserService.UserRegister(ctx, req)
	if err != nil {
		return
	}

	return
}
