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
	ListUsers(ctx context.Context) ([]*model.UserInfo, error)
	UpdateUser(ctx context.Context, id uint64, userMap map[string]interface{}) error
	AdminUpdateUser(ctx context.Context, id uint64, username, email, password, role string) error
	DeleteUser(ctx context.Context, id uint64) error
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

	payload := map[string]interface{}{"user_id": u.ID, "username": u.Username, "role": u.Role}
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

// ListUsers 列出所有用户（管理员用）
func (s *userService) ListUsers(ctx context.Context) ([]*model.UserInfo, error) {
	users, err := s.repo.ListUsers(ctx)
	if err != nil {
		return nil, err
	}
	infos := make([]*model.UserInfo, 0, len(users))
	for _, u := range users {
		infos = append(infos, u.ToUserInfo())
	}
	return infos, nil
}

// AdminUpdateUser 管理员修改用户信息（除头像外的字段）
// 传入空字符串表示不修改该字段
func (s *userService) AdminUpdateUser(ctx context.Context, id uint64, username, email, password, role string) error {
	userMap := make(map[string]interface{})
	if username != "" {
		userMap["username"] = username
	}
	if email != "" {
		userMap["email"] = email
	}
	if role != "" {
		if role != model.RoleAdmin && role != model.RoleUser {
			return errors.New("invalid role")
		}
		userMap["role"] = role
	}
	if password != "" {
		pwd, err := auth.HashAndSalt(password)
		if err != nil {
			return pkgerrors.Wrap(err, "encrypt password err")
		}
		userMap["password"] = pwd
	}
	if len(userMap) == 0 {
		return nil
	}
	userMap["updated_at"] = time.Now()
	return s.repo.UpdateUser(ctx, id, userMap)
}

// DeleteUser 管理员删除用户
func (s *userService) DeleteUser(ctx context.Context, id uint64) error {
	return s.repo.DeleteUser(ctx, id)
}
