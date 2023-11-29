package service

import (
	"context"
	"errors"
	"sync"

	"gorm.io/gorm"

	"github.com/CocaineCong/micro-todoList/app/user/repository/db/dao"
	"github.com/CocaineCong/micro-todoList/app/user/repository/db/model"
	"github.com/CocaineCong/micro-todoList/idl/pb"
	"github.com/CocaineCong/micro-todoList/pkg/e"
)

var UserSrvIns *UserSrv
var UserSrvOnce sync.Once

type UserSrv struct {
}

func GetUserSrv() *UserSrv {
	UserSrvOnce.Do(func() {
		UserSrvIns = &UserSrv{}
	})
	return UserSrvIns
}

func (u *UserSrv) UserLogin(ctx context.Context, req *pb.UserRequest, resp *pb.UserDetailResponse) (err error) {
	resp.Code = e.SUCCESS
	user, err := dao.NewUserDao(ctx).FindUserByUserName(req.UserName)
	if err != nil {
		resp.Code = e.ERROR
		return
	}

	if !user.CheckPassword(req.Password) {
		resp.Code = e.InvalidParams
		return
	}

	resp.UserDetail = BuildUser(user)
	return
}

func (u *UserSrv) UserRegister(ctx context.Context, req *pb.UserRequest, resp *pb.UserDetailResponse) (err error) {
	if req.Password != req.PasswordConfirm {
		err = errors.New("两次密码输入不一致")
		return
	}
	resp.Code = e.SUCCESS
	_, err = dao.NewUserDao(ctx).FindUserByUserName(req.UserName)
	if err != nil {
		if err == gorm.ErrRecordNotFound { // 如果不存在就继续下去
			// ...continue
		} else {
			resp.Code = e.ERROR
			return
		}
	}
	user := &model.User{
		UserName: req.UserName,
	}
	// 加密密码
	if err = user.SetPassword(req.Password); err != nil {
		resp.Code = e.ERROR
		return
	}
	if err = dao.NewUserDao(ctx).CreateUser(user); err != nil {
		resp.Code = e.ERROR
		return
	}

	resp.UserDetail = BuildUser(user)
	return
}

func BuildUser(item *model.User) *pb.UserModel {
	userModel := pb.UserModel{
		Id:        uint32(item.ID),
		UserName:  item.UserName,
		CreatedAt: item.CreatedAt.Unix(),
		UpdatedAt: item.UpdatedAt.Unix(),
	}
	return &userModel
}
