package pgx

import "errors"

var (
	ErrUserExists = errors.New("user already exists")
	ErrUserCreate = errors.New("failed to create user")
)
