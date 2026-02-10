package errs

import (
	"encoding/json"
	"fmt"
)

type Code uint

const (
	CodeNotFound Code = iota
	CodeAlreadyExists
	CodeBadRequest
	CodeInternal
)

type Error struct {
	Code    Code   `json:"-"`
	Message string `json:"message"`
	Cause   error  `json:"-"`
}

func NewError(code Code, message string) *Error {
	return &Error{Code: code, Message: message}
}

func (e *Error) Error() string {
	data, err := json.Marshal(e)
	if err != nil {
		return "invalid error"
	}
	return string(data)
}

func (e *Error) WithCause(err error) *Error {
	e.Cause = err
	return e
}

func (e *Error) SetCause(err error) {
	e.Cause = err
}

func (e *Error) SetCode(code Code) {
	e.Code = code
}

func (e *Error) SetMessage(message string) {
	e.Message = message
}

func NewEntityNotFoundError(entityName string) *Error {
	return NewError(CodeNotFound, fmt.Sprintf("%s not found", entityName))
}

func NewEntityAlreadyExistsError(entityName string) *Error {
	return NewError(CodeAlreadyExists, fmt.Sprintf("%s already exists", entityName))
}

func NewBadRequestError(message string) *Error {
	return NewError(CodeBadRequest, message)
}

func NewInternalServerError() *Error {
	return NewError(CodeInternal, "Internal server error")
}
