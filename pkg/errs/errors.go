package errs

import "errors"

var (
	ErrAlreadyExists = errors.New("this user already exists")
	ErrNotFound      = errors.New("not found")
	ErrInvalidCreds  = errors.New("invalid email or password")
)
