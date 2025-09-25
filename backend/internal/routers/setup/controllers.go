package setup

import (
	"insight/internal/controller/demo"
	"insight/internal/controller/hello"
)

type Controllers struct {
	HelloController hello.HelloController
	DemoController  demo.DemoController
}

// NewControllers creates and returns a Controllers instance with its HelloController
// field initialized.
//
// It instantiates a hello.HelloController via hello.NewHelloController and returns
// a pointer to the populated Controllers struct.
func NewControllers() *Controllers {

	HelloController := hello.NewHelloController()
	DemoController := demo.NewDemoController()
	return &Controllers{
		HelloController: *HelloController,
		DemoController:  *DemoController,
	}
}
