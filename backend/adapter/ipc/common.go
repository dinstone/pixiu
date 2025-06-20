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

func Failure(err error) *Result {
	result := &Result{}
	if err != nil {
		handleError(err, result)
	} else {
		result.Code = 599
	}
	return result
}

func Success(data any) *Result {
	return &Result{
		Code: 0,
		Mesg: "ok",
		Data: data,
	}
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

func handleSuccess(token any, result *Result) {
	result.Data = token
}
