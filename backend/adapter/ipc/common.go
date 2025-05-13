package ipc

import (
	"errors"
	"pixiu/backend/pkg/exception"
)

type Result struct {
	Code int    `json:"code"`
	Mesg string `json:"mesg"`
	Data any    `json:"data,omitempty"`
}

func handleError(err error, result *Result) {
	var appErr exception.AppError
	if errors.As(err, &appErr) {
		result.Code = appErr.Code()
		result.Mesg = appErr.Error()
	} else {
		result.Code = 500
		result.Mesg = err.Error()
	}
}
