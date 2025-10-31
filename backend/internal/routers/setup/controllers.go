package setup

import (
	"insight/internal/controller/admin"
	"insight/internal/controller/demo"
	"insight/internal/controller/hello"
)

type Controllers struct {
	HelloController      hello.HelloController
	DemoController       demo.DemoController
	UserController       admin.AdminUserController
	LoginController      admin.LoginController
	PermissionController admin.PermissionController
	RoleController       admin.RoleController
}

// NewControllers creates and returns a Controllers instance with its HelloController
// field initialized.
//
// It instantiates a hello.HelloController via hello.NewHelloController and returns
// a pointer to the populated Controllers struct.
func NewControllers() *Controllers {

	HelloController := hello.NewHelloController()
	DemoController := demo.NewDemoController()
	UserController := admin.NewAdminUserController()
	LoginController := admin.NewLoginController()
	PermissionController := admin.NewPermissionController()
	RoleController := admin.NewRoleController()

	return &Controllers{
		HelloController:      *HelloController,
		DemoController:       *DemoController,
		UserController:       *UserController,
		LoginController:      *LoginController,
		PermissionController: *PermissionController,
		RoleController:       *RoleController,
	}
}
