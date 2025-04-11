package errors

import "errors"

type MapError struct {
	Code        int    `json:"-"`
	Message     string `json:"message"`
	Cause       error  `json:"-"`
	Description string `json:"description,omitempty"`
}

func (e *MapError) Error() string {
	if e.Cause != nil {
		return e.Message + ": " + e.Cause.Error()
	}
	return e.Message
}

func (e *MapError) Is(target error) bool {
	if err, ok := target.(*MapError); ok {
		return err.Code == e.Code && err.Message == e.Message
	}
	return errors.Is(e.Cause, target)
}

func NewError(code int, message string) *MapError {
	return &MapError{Code: code, Message: message}
}

func UnWrapError(code int, message string, description string) *MapError {
	return &MapError{Code: code, Message: message, Description: description}
}

func WrapErrorWithCause(code int, message string, cause error) *MapError {
	return &MapError{Code: code, Message: message, Cause: cause}
}
