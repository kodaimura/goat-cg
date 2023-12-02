package errs

import (
	"fmt"
)


type UniqueConstraintError struct {
	Column string
}

func NewUniqueConstraintError(column string) error {
	return UniqueConstraintError{Column: column}
}

func (e UniqueConstraintError) Error() string {
	return fmt.Sprintf("UniqueConstraintError: %s", e.Column)
}