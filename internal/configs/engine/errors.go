package engine

import (
	"fmt"
	"net/http"
)

type Error struct {
	Code  int                    `json:"code"`
	Meta  MetaData               `json:"meta"`
	Extra map[string]interface{} `json:"extra"`
}

func NewGenericError(code int, message string) *Error {
	return &Error{Meta: MetaData{Type: "ERROR", Message: message}, Extra: make(map[string]interface{}), Code: code}
}

func ErrBadRequest() *Error {
	return NewGenericError(http.StatusBadRequest, "Bad Request")
}

func ErrNotFound() *Error {
	return NewGenericError(http.StatusNotFound, "Not Found")
}

func ErrPrecondRequired() *Error {
	return NewGenericError(http.StatusPreconditionRequired, "Precondition Required")
}

func ErrUnauthorized() *Error {
	return NewGenericError(http.StatusUnauthorized, "Unauthorized Access")
}

func ErrInvalidPassword() *Error {
	return NewGenericError(http.StatusUnauthorized, "Invalid Password")
}

func ErrInternalFailure() *Error {
	return NewGenericError(http.StatusInternalServerError, "Internal Server Error")
}

func (e *Error) Message(message string) *Error {
	e.Meta.Message = message
	return e
}

func (e *Error) ExtraData(data map[string]interface{}) *Error {
	e.Extra = data
	return e
}

func (e Error) Error() string {
	return fmt.Sprintf("error [%d]: %s", e.Code, e.Meta.Message)
}

func (e Error) SimpleError() error {
	return fmt.Errorf("error [%d]: %s", e.Code, e.Meta.Message)
}
