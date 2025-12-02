package resources

import (
	"errors"

	"github.com/jinzhu/copier"
)

type AdminUserResources struct {
	ID       uint     `json:"userId"`
	IsAdmin  int8     `json:"isAdmin"`
	NickName string   `json:"realName"`
	Username string   `json:"username"`
	Email    string   `json:"email"`
	Avatar   string   `json:"avatar"`
	Mobile   string   `json:"mobile"`
	Roles    []string `json:"roles"`
	homePath string   `json:"homePath"`
	token    string   `json:"token"`
}

func NewAdminUserResources(data any) *AdminUserResources {
	var adminUser AdminUserResources
	err := copier.Copy(&adminUser, data)
	if err != nil {
		return nil
	}
	return &adminUser
}

func (r *AdminUserResources) SetToken(token string) error {
	if token == "" {
		return errors.New("token is empty")
	}
	r.token = token
	return nil
}

func (r *AdminUserResources) SetRoles(roles []string) {
	r.Roles = roles
}
