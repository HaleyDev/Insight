package resources

import (
	"errors"
	"github.com/jinzhu/copier"
	"insight/internal/model"
)

type AdminUserResources struct {
	avatar   string   `json:"avatar"`
	realName string   `json:"realName"`
	roles    []string `json:"roles"`
	userId   string   `json:"userId"`
	username string   `json:"username"`
	desc     string   `json:"desc"`
	homePath string   `json:"homePath"`
	token    string   `json:"token"`
}

func NewAdminUserResources(data model.AdminUser) *AdminUserResources {
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
