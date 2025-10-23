package validator

import (
	"reflect"
	"regexp"
	"sync"

	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

var trans ut.Translator  // 全局校验器

var once sync.Once

var validate *validator.Validate

func InitValidatorTrans(locale string) {
	once.Do(func() { v })
}

func validatorTrans(locale string) {
	var ok bool
	if validate, ok = binding.Validator.Engine().(*validator.Validate); !ok {
		panic("Failed to initialize the validator")
	}
	// 注册自定义验证消息
	registerValidation()
	// 注册一个获取json tag的自定义方法
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {


}


func registerValidation() {
	// 注册手机号验证规则
	err := validate.RegisterValidation("mobile", func(fl validator.FieldLevel) bool {
		return regexp.MustCompile(`^1[3456789]\d{9}$`).MatchString(fl.Field().String())
	})
	if err != nil {
		panic("Failed to register mobile validation")
	}
}