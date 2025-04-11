package errors

import "errors"

// MapError represents API error response
// @swagger:model
type MapError struct {
	// HTTP status code
	// example: 400
	Code int `json:"-"`

	// Error message
	// example: "invalid request parameters"
	Message string `json:"message"`

	// Internal error cause (not exposed to clients)
	// example: "invalid request parameters"
	Cause error `json:"-"`

	// Detailed error description
	// example: "Field 'email' must be valid email address"
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
