package errs

import "encoding/json"

type Code uint

const (
	CodeUnautorized Code = iota
	CodeForbidden
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

func NewUnauthorizedError() *Error {
	return NewError(CodeUnautorized, "The token is missing, invalid or expired")
}

func NewForbiddenError() *Error {
	return NewError(CodeForbidden, "The token is valid, but lacks permission")
}

func NewInternalServerError() *Error {
	return NewError(CodeInternal, "Internal server error")
}
