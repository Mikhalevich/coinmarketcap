package coinmarketcap

import (
	"fmt"
)

type Error struct {
	Code    int
	Message string
}

func NewError(code int, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

func (e *Error) Error() string {
	return fmt.Sprintf("code: %d message: %s", e.Code, e.Message)
}
