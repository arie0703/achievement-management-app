package errors

import (
	"errors"
	"fmt"
)

// Common error types
var (
	ErrNotFound           = errors.New("resource not found")
	ErrInvalidInput       = errors.New("invalid input")
	ErrInsufficientPoints = errors.New("insufficient points")
	ErrDuplicateResource  = errors.New("resource already exists")
	ErrDatabaseOperation  = errors.New("database operation failed")
)

// ValidationError バリデーションエラー
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation error for field '%s': %s", e.Field, e.Message)
}

// BusinessLogicError ビジネスロジックエラー
type BusinessLogicError struct {
	Operation string
	Reason    string
}

func (e BusinessLogicError) Error() string {
	return fmt.Sprintf("business logic error in operation '%s': %s", e.Operation, e.Reason)
}

// DatabaseError データベースエラー
type DatabaseError struct {
	Operation string
	Table     string
	Cause     error
}

func (e DatabaseError) Error() string {
	return fmt.Sprintf("database error in operation '%s' on table '%s': %v", e.Operation, e.Table, e.Cause)
}

func (e DatabaseError) Unwrap() error {
	return e.Cause
}

// ServiceError サービス層エラー
type ServiceError struct {
	Operation string
	Message   string
	Cause     error
}

func (e ServiceError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("service error in operation '%s': %s (caused by: %v)", e.Operation, e.Message, e.Cause)
	}
	return fmt.Sprintf("service error in operation '%s': %s", e.Operation, e.Message)
}

func (e ServiceError) Unwrap() error {
	return e.Cause
}