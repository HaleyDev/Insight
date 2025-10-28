package validator

import (
	"reflect"
	"regexp"
	"sync"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
)

var trans ut.Translator // 全局校验器

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
		label := field.Tag.Get("label")
		if label == "" {
			label = field.Tag.Get("json")
			if label == "" {
				label = field.Tag.Get("form")
			}
		}

		if label == "-" {
			return ""
		}

		if label == "" {
			return field.Name
		}

		return label
	})

	zhT := zh.New() // 中文翻译器
	enT := en.New() // 英文翻译器
	uni := ut.New(enT, zhT, enT)

	trans, ok = uni.GetTranslator(locale)
	if !ok {
		panic("Initialize a language not supported by the validator")
	}

	var err error

	switch locale {
	case "en":
		err = enTranslations.RegisterDefaultTranslations(validate, trans)
	case "zh":
		err = zhTranslations.RegisterDefaultTranslations(validate, trans)
	default:
		err = enTranslations.RegisterDefaultTranslations(validate, trans)
	}
	if err != nil {
		panic("Failed to register translator when initializing validator")
	}

	err = 
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


type translation struct {
	tag string
	translation string
	override bool
	customRegisFunc validator.RegisterTranslationsFunc
	customTransFunc validator.TranslationFunc
}

func customRegisTranslation() error {
	translations := []translation{
		{
			tag: "mobile",
			translation: "{0}格式不正确",
			override: false,
		},
	}
}


func registerTranslation(translations []translation) (err error) {
	for _, t := range translations {
		if t.customTransFunc != nil && t.customRegisFunc != nil {
			err = validate.RegisterTranslation(t.tag, trans, t.customRegisFunc, t.customTransFunc)
		} else if t.customTransFunc != nil && t.customRegisFunc == nil {
			err = validate.RegisterTranslation(t.tag, trans, registrationFunc(t.tag, t.translation, t.override), t.customTransFunc)
		} else if t.customTransFunc == nil && t.customRegisFunc != nil {
			err = validate.RegisterTranslation(t.tag, trans, t.customRegisFunc, translateFunc)
		} else {
			err = validate.RegisterTranslation(t.tag, trans, registrationFunc(t.tag, t.translation, t.override), translateFunc)
		}

		if err != nil {
			return
		}
	}
	return
}
