package domain

import (
	"github.com/pkg/errors"
)

type ErrorType int

func (t ErrorType) Code() int {
	return int(t)
}

const (
	// 内部エラー
	ErrorTypeInternal ErrorType = 500
	// クライアントリクエスト起因エラー
	ErrorTypeBadRequest = 400
	// 認証エラー
	ErrorTypeUnauthorized = 401
	// 権限エラー
	ErrorTypeForbidden = 403
	// 存在しないリソース参照エラー
	ErrorTypeNotFound = 404
	// コンフリクトエラー
	ErrorTypeAlreadyExists = 409
)

const (
	ForbiddenAgencyMsg     = "不正な代理店です"
	ForbiddenDepartmentMsg = "不正な店舗です"
	ConflictStatusMsg      = "現在のステータスから要求されたステータスへの変更は許可されていません"
	BadRequestDateMsg      = "指定された日付が不正です"
	BadRequestMsg          = "パラメータが不正です"
)

func NewBadRequestErr(message string) error {
	return NewError(ErrorTypeBadRequest, message)
}

func NewUnAuthorizedErr(message string) error {
	return NewError(ErrorTypeUnauthorized, message)
}

func NewForbiddenErr(message string) error {
	return NewError(ErrorTypeForbidden, message)
}

func NewNotFoundErr() error {
	return NewError(ErrorTypeNotFound, "指定されたリソースは存在しません")
}

func NewConflictErr(message string) error {
	return NewError(ErrorTypeAlreadyExists, message)
}

func NewInternalServerErr() error {
	return NewError(ErrorTypeInternal, "内部エラーです")
}

type AppError interface {
	error
	Type() ErrorType
}

type appError struct {
	errType ErrorType
	message string
}

func NewError(errType ErrorType, message string) error {
	return &appError{
		errType: errType,
		message: message,
	}
}

func (err *appError) Error() string {
	return err.message
}

func (err *appError) Type() ErrorType {
	return err.errType
}

func IsNotFound(err error) bool {
	appErr, ok := errors.Cause(err).(AppError)
	if !ok {
		return false
	}

	return appErr.Type() == ErrorTypeNotFound
}

func IsConflict(err error) bool {
	appErr, ok := errors.Cause(err).(AppError)
	if !ok {
		return false
	}

	return appErr.Type() == ErrorTypeAlreadyExists
}
