package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/CocaineCong/micro-todoList/app/gateway/rpc"
	"github.com/CocaineCong/micro-todoList/idl/pb"
	"github.com/CocaineCong/micro-todoList/pkg/ctl"
	log "github.com/CocaineCong/micro-todoList/pkg/logger"
	"github.com/CocaineCong/micro-todoList/pkg/utils"
	"github.com/CocaineCong/micro-todoList/types"
)

// UserRegisterHandler 用户注册处理函数
func UserRegisterHandler(ctx *gin.Context) {
	var req pb.UserRequest
	if err := ctx.ShouldBind(&req); err != nil {
		// 如果参数绑定失败，返回400错误
		ctx.JSON(http.StatusBadRequest, ctl.RespError(ctx, err, "UserRegister Bind 绑定参数失败"))
		return
	}
	// 参数校验成功后，发起RPC调用
	userResp, err := rpc.UserRegister(ctx, &req)
	if err != nil {
		// 如果RPC调用失败，记录错误并返回500错误
		log.LogrusObj.Errorf("UserRegister:%v", err)
		ctx.JSON(http.StatusInternalServerError, ctl.RespError(ctx, err, "UserRegister RPC 调用失败"))
		return
	}
	// 如果成功，返回200和用户响应
	ctx.JSON(http.StatusOK, ctl.RespSuccess(ctx, userResp))
}

// UserLoginHandler 用户登录处理函数
func UserLoginHandler(ctx *gin.Context) {
	var req pb.UserRequest
	if err := ctx.Bind(&req); err != nil {
		// 如果参数绑定失败，返回400错误
		ctx.JSON(http.StatusBadRequest, ctl.RespError(ctx, err, "UserLogin Bind 绑定参数失败"))
		return
	}
	// 参数校验成功后，发起RPC调用
	userResp, err := rpc.UserLogin(ctx, &req)
	if err != nil {
		// 如果RPC调用失败，返回500错误
		ctx.JSON(http.StatusInternalServerError, ctl.RespError(ctx, err, "UserLogin RPC 调用失败"))
		return
	}
	// 生成令牌
	token, err := utils.GenerateToken(uint(userResp.UserDetail.Id))
	if err != nil {
		// 如果令牌生成失败，返回500错误
		ctx.JSON(http.StatusInternalServerError, ctl.RespError(ctx, err, "GenerateToken 失败"))
		return
	}
	// 封装响应体，包括用户信息和令牌
	res := &types.TokenData{
		User:  userResp,
		Token: token,
	}
	// 如果成功，返回200和响应体
	ctx.JSON(http.StatusOK, ctl.RespSuccess(ctx, res))
}
