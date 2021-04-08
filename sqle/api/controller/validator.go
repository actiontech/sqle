package controller

import (
	"fmt"
	"github.com/go-playground/validator"
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
)

var (
	ValidateNameRegexpPattern = `^[a-zA-Z][a-zA-Z0-9\_\-]{0,59}$`
)

type CustomValidator struct {
	validate *validator.Validate
	enTrans  ut.Translator
	zhTrans  ut.Translator
}

var DefaultCustomValidator *CustomValidator

func init() {
	cv := NewCustomValidator()
	err := cv.RegisterTranslation("required", "{0} is required", "{0}不能为空")
	if err != nil {
		panic(err)
	}
	err = cv.RegisterTranslation("oneof", "{0} must be one of [{1}]", "{0}必须是[{1}]其中之一")
	if err != nil {
		panic(err)
	}
	err = cv.RegisterTranslation("email", "{0} is invalid email addr", "{0}是无效的邮箱地址")
	if err != nil {
		panic(err)
	}
	err = cv.RegisterTranslation("name", "{0} must match regexp `{1}`", "{0}必须匹配正则`{1}`",
		ValidateNameRegexpPattern)
	if err != nil {
		panic(err)
	}
	DefaultCustomValidator = cv
}

func NewCustomValidator() *CustomValidator {
	en := en.New()
	zh := zh.New()
	uni := ut.New(en, zh)

	cv := &CustomValidator{
		validate: validator.New(),
	}
	enTrans, _ := uni.GetTranslator("en")
	zhTrans, _ := uni.GetTranslator("zh")
	cv.enTrans = enTrans
	cv.zhTrans = zhTrans

	cv.validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	cv.validate.SetTagName("valid")

	cv.validate.RegisterValidation("name", ValidateName)
	return cv
}

func (cv *CustomValidator) RegisterTranslation(tag, enText, zhText string, params ...string) error {
	err := cv.validate.RegisterTranslation(tag, cv.enTrans, func(ut ut.Translator) error {
		return ut.Add(tag, enText, true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		tParams := []string{fe.Namespace()}
		if len(params) > 0 {
			tParams = append(tParams, params...)
		}
		tParams = append(tParams, fe.Param())
		t, _ := ut.T(tag, tParams...)
		return t
	})
	if err != nil {
		return err
	}
	err = cv.validate.RegisterTranslation(tag, cv.zhTrans, func(ut ut.Translator) error {
		return ut.Add(tag, zhText, true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		tParams := []string{fe.Namespace()}
		if len(params) > 0 {
			tParams = append(tParams, params...)
		}
		tParams = append(tParams, fe.Param())
		t, _ := ut.T(tag, tParams...)
		return t
	})
	return err
}

func (cv *CustomValidator) Validate(i interface{}) error {
	err := cv.validate.Struct(i)
	if err == nil {
		return nil
	}
	errs := err.(validator.ValidationErrors)
	errMsgs := make([]string, 0, len(errs))
	for _, err := range errs {
		errMsgs = append(errMsgs, err.Translate(cv.enTrans))
	}
	if len(errMsgs) > 0 {
		return fmt.Errorf(strings.Join(errMsgs, "; "))
	}
	return nil
}

// ValidateMyVal implements validator.Func
func ValidateName(fl validator.FieldLevel) bool {
	return validateName(fl.Field().String())
}

func validateName(name string) bool {
	match, _ := regexp.MatchString(ValidateNameRegexpPattern, name)

	return match
}
