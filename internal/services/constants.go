package services

import "errors"

const (
	PgErrCodeUniqueViolation = "23505"
)

var (
	ErrDuplicatedUsernameOrEmail = errors.New("username or email already exists")
	ErrInvalidCredentials        = errors.New("invalid credentials")
)
