package errors

import "fmt"

type AppError struct {
	StatusCode int
	Message    string
}

func (e *AppError) Error() string {
	return e.Message
}

func New(statusCode int, message string) *AppError {
	return &AppError{StatusCode: statusCode, Message: message}
}

func BadRequest(message string) *AppError {
	return New(400, message)
}

func Unauthorized(message string) *AppError {
	return New(401, message)
}

func Forbidden(message string) *AppError {
	return New(403, message)
}

func NotFound(message string) *AppError {
	return New(404, message)
}

func Conflict(message string) *AppError {
	return New(409, message)
}

func Internal(message string) *AppError {
	return New(500, message)
}

func Wrap(statusCode int, format string, args ...any) *AppError {
	return New(statusCode, fmt.Sprintf(format, args...))
}
