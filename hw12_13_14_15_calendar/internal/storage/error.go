package storage

import "github.com/pkg/errors"

var (
	ErrDateBusy = errors.New("date is busy")
	ErrNotFound = errors.New("event not found")
)
