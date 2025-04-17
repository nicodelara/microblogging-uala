package application

import "errors"

var (
	ErrUserNotFound   = errors.New("user does not exist")
	ErrContentTooLong = errors.New("content exceeds 280 characters")
)
