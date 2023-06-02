package domain

import "errors"

var (
	ErrUserNotFound      = errors.New("user doesn't exists")
	ErrUserAlreadyExists = errors.New("user with such email or username is already exists")
)
