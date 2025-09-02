package repository

import "errors"

var (
	NotFound         = errors.New("uuid not found")
	InternalError    = errors.New("failed to complete task")
	ErrOrderExists   = errors.New("order already exists")
	ErrOrderNotFound = errors.New("order id is incorrect")
)
