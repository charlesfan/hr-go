package db

import (
	"errors"
	"strings"
)

var (
	duplicateMsg = "ERROR: duplicate key value violates unique constraint"
)

var (
	ErrDuplicate     = errors.New("Duplicate column error")
	ErrDuplicateName = errors.New("Duplicate name error")
	ErrNotFound      = errors.New("Not found")
)

func ErrChecking(err error) error {
	if strings.Contains(err.Error(), duplicateMsg) {
		switch {
		case strings.Contains(err.Error(), "name_key"):
			return ErrDuplicateName
		case strings.Contains(err.Error(), "not found"):
			return ErrNotFound
		default:
			return ErrDuplicate
		}
	}
	return err
}
