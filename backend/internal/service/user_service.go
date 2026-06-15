package service

import (
	"context"
	"errors"
	"time"

	pkgerrors "github.com/pkg/errors"

	"github.com/insight/backend/internal/model"
	"github.com/insight/backend/internal/repository"
	"github.com/insight/backend/pkg/app"
	"github.com/insight/backend/pkg/auth"
)

// UserService 用户服务接口
type UserService interface {
	Register(ctx context.Context, username, email, password string) error
	EmailLogin(ctx context.Context, email, password string) (*model.LoginResult, error)
	GetUserByID(ctx context.Context, id uint64) (*model.UserBaseModel, error)
	GetUserInfoByID(ctx context.Context, id uint64) (*model.UserInfo, error)
	GetUserByEmail(ctx context.Context, email string) (*model.UserBaseModel, error)
	UpdateUser(ctx context.Context, id uint64, userMap map[string]interface{}) error
}

type userService struct {
	repo repository.Repository
}

var _ UserService = (*userService)(nil)

func newUsers(svc *service) *userService {
	return &userService{repo: svc.repo}
}

// Register 注册用户：默认头像 + 默认角色 user
func (s *userService) Register(ctx context.Context, username, email, password string) error {
	pwd, err := auth.HashAndSalt(password)
	if err != nil {
		return pkgerrors.Wrapf(err, "encrypt password err")
	}

	u := &model.UserBaseModel{
		Username:  username,
		Password:  pwd,
		Email:     email,
		Avatar:    model.DefaultAvatar,
		Role:      model.RoleUser,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	isExist, err := s.repo.UserIsExist(u)
	if err != nil {
		return pkgerrors.Wrapf(err, "check user exist err")
	}
	if isExist {
		return errors.New("用户名或邮箱已被注册")
	}

	if _, err = s.repo.CreateUser(ctx, u); err != nil {
		return pkgerrors.Wrapf(err, "create user err")
	}
	return nil
}

// EmailLogin 邮箱登录，返回 token + 用户信息
func (s *userService) EmailLogin(ctx context.Context, email, password string) (*model.LoginResult, error) {
	u, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, pkgerrors.Wrap(err, "get user info err by email")
	}

	if !auth.ComparePasswords(u.Password, password) {
		return nil, errors.New("invalid password")
	}

	payload := map[string]interface{}{"user_id": u.ID, "username": u.Username}
	tokenStr, err := app.Sign(ctx, payload, app.Conf.JwtSecret, 86400)
	if err != nil {
		return nil, pkgerrors.Wrap(err, "gen token err")
	}

	return &model.LoginResult{
		Token: tokenStr,
		User:  u.ToUserInfo(),
	}, nil
}

// GetUserByID 获取单条用户记录
func (s *userService) GetUserByID(ctx context.Context, id uint64) (*model.UserBaseModel, error) {
	return s.repo.GetUser(ctx, id)
}

// GetUserInfoByID 获取对外用户结构
func (s *userService) GetUserInfoByID(ctx context.Context, id uint64) (*model.UserInfo, error) {
	u, err := s.repo.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}
	return u.ToUserInfo(), nil
}

// GetUserByEmail 根据邮箱获取
func (s *userService) GetUserByEmail(ctx context.Context, email string) (*model.UserBaseModel, error) {
	return s.repo.GetUserByEmail(ctx, email)
}

// UpdateUser 更新用户字段
func (s *userService) UpdateUser(ctx context.Context, id uint64, userMap map[string]interface{}) error {
	return s.repo.UpdateUser(ctx, id, userMap)
}
