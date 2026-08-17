package util

import (
	"errors"
	"fmt"
	"net/http"
)

// AppError 统一业务错误，携带业务 code 与 HTTP 状态码。
type AppError struct {
	Code       int
	HTTPStatus int
	Message    string
	Err        error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("code=%d message=%s cause=%v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("code=%d message=%s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error { return e.Err }

// NewAppError 构造业务错误。
func NewAppError(code, httpStatus int, message string) *AppError {
	return &AppError{Code: code, HTTPStatus: httpStatus, Message: message}
}

// WrapAppError 包裹底层错误并保留业务错误信息。
func WrapAppError(appErr *AppError, cause error) *AppError {
	appErr.Err = cause
	return appErr
}

// IsAppError 判断错误是否为 AppError。
func IsAppError(err error) bool {
	var appErr *AppError
	return errors.As(err, &appErr)
}

// ToAppError 将任意错误转为 AppError，未知错误统一为 500。
func ToAppError(err error) *AppError {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}
	return WrapAppError(NewAppError(50000, http.StatusInternalServerError, "系统繁忙，请稍后重试"), err)
}
