package validator

import (
	"github.com/go-playground/validator/v10"

	"git.ocn.com.vn/ocn/common/httpbase"
	"git.ocn.com.vn/ocn/common/logger"
)

func registerValidator(tagName string, validatorFunc validator.Func) {
	if err := httpbase.RegisterValidator(tagName, validatorFunc); err != nil {
		panic(err)
	} else {
		logger.Debug("register validator ", "tag name", tagName)
	}
}

func registerStructValidator(fn validator.StructLevelFunc, in ...interface{}) {
	httpbase.RegisterStructValidation(fn, in...)
}

func Init() {
	// register tag or struct here
}
