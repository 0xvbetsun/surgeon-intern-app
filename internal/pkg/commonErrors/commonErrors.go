package commonErrors

import "fmt"

type (
	NotExistsError struct {
		Resource string
		ID       string
	}
	InvalidInputError struct {
		message string
		cause   error
	}
	UnauthorizedError struct {
		message string
		cause   error
	}
	NotImplementedError struct {
		message string
	}
)

func NewNotImplementedError() NotImplementedError {
	return NotImplementedError{message: "Method not implemented"}
}

func NewUnauthorizedAccessError(message string) UnauthorizedError {
	if message != "" {
		return UnauthorizedError{
			message: message,
		}
	}
	return UnauthorizedError{
		message: "Access Denied",
	}
}

func NewUnauthorizedAccessErrorWithCause(message string, cause error) UnauthorizedError {
	if message != "" {
		return UnauthorizedError{
			message: message,
			cause:   cause,
		}
	}
	return UnauthorizedError{
		message: "Access Denied",
		cause:   cause,
	}
}

func (n NotImplementedError) Error() string {
	return n.message
}

func (u UnauthorizedError) Error() string {
	return u.message
}

func (u UnauthorizedError) Cause() error {
	return u.cause
}

func (n NotExistsError) Error() string {
	return fmt.Sprintf("The requested %s %s does not exist", n.Resource, n.ID)
}

func NewInvalidInputError(message string) InvalidInputError {
	return InvalidInputError{
		message: message,
	}
}

func NewInvalidInputErrorWithCause(message string, cause error) InvalidInputError {
	return InvalidInputError{
		message: message,
		cause:   cause,
	}
}

func (e InvalidInputError) Error() string {
	return e.message
}

func (e InvalidInputError) Cause() error {
	return e.cause
}
