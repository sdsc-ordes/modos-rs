package types

import "fmt"

var _ error = (*err)(nil)
var _ error = (*PermissionError)(nil)
var _ error = (*DurationError)(nil)

type (
	err struct {
		Message string
	}

	PermissionError struct {
		err
	}

	DurationError struct {
		err
	}
)

func (err *err) Error() string {
	return err.Message
}

func NewPermissionsError(msg string, args ...any) error {
	return &PermissionError{err{Message: fmt.Sprintf(msg, args...)}}
}

func NewDurationError(msg string, args ...any) error {
	return &DurationError{err{Message: fmt.Sprintf(msg, args...)}}
}
