package domain

import "errors"

var (
	ErrUserAlreadyExists  = errors.New("username already exists")
	ErrNotFound           = errors.New("resource not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
)
