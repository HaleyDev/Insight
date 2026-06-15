package ecode

import "github.com/insight/backend/pkg/errcode"

//nolint: golint
var (
	ErrUserNotFound          = errcode.NewError(20101, "用户不存在")
	ErrPasswordIncorrect     = errcode.NewError(20102, "账号或密码错误")
	ErrEmailOrPassword       = errcode.NewError(20109, "邮箱或密码错误")
	ErrTwicePasswordNotMatch = errcode.NewError(20110, "两次密码输入不一致")
	ErrRegisterFailed        = errcode.NewError(20111, "注册失败")
)
