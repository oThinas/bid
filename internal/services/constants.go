package services

import (
	"errors"
	"time"
)

const (
	PgErrCodeUniqueViolation = "23505"
	MaxMessageSize           = 512
	ReadDeadline             = 60 * time.Second
	WriteDeadLine            = 10 * time.Second
	PingInterval             = (ReadDeadline * 9) / 10
)

var (
	ErrDuplicatedUsernameOrEmail = errors.New("username or email already exists")
	ErrInvalidCredentials        = errors.New("invalid credentials")
	ErrBidAmountTooLow           = errors.New("bid amount is too low")
	ErrProductNotFound           = errors.New("product not found")
)
