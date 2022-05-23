package repoerrors

import (
	"fmt"
	"github.com/friendsofgo/errors"
	"github.com/lib/pq"
)

type (
	BadArgumentError struct {
		Message string
	}
	NotFoundError struct {
		message string
		cause   error
	}
	QueryError struct {
		Message string
		cause   error
	}
)

func NewNotFoundError(entityID interface{}) NotFoundError {
	return NotFoundError{
		message: fmt.Sprintf("Entity with id %v not found", entityID),
	}
}

func NewNotFoundErrorWithCause(entityID interface{}, cause error) NotFoundError {
	return NotFoundError{
		message: fmt.Sprintf("Entity with id %v not found", entityID),
		cause:   cause,
	}
}

func (e NotFoundError) Error() string {
	return e.message
}

func (e NotFoundError) Cause() error {
	return e.cause
}

func NewBadArgumentError(argumentName string) BadArgumentError {
	return BadArgumentError{
		Message: fmt.Sprintf("Illegal argument: %s", argumentName),
	}
}

func (e BadArgumentError) Error() string {
	return e.Message
}

func NewQueryError(err *pq.Error) QueryError {
	return QueryError{
		Message: fmt.Sprintf("Query failed PG Error Code Id: %s Name: %s Message: %s Details: %s", err.Code, err.Code.Name(), err.Message, err.Detail),
		cause:   err,
	}
}

func (e QueryError) Error() string {
	return e.Message
}

func (e QueryError) Cause() error {
	return e.cause
}

func ErrorFromDbError(err error) error {
	if driverErr, ok := err.(*pq.Error); ok {
		return NewQueryError(driverErr)
	}
	if driverErr, ok := errors.Cause(err).(*pq.Error); ok {
		return NewQueryError(driverErr)
	}
	return errors.WithMessage(err, "Unknown Database Error")
}
