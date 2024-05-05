package app_errors

import "errors"

var (
	ErrUserAlreadyExists = errors.New("user is already exist")
	ErrInvalidPassword   = errors.New("password is invalid")

	ErrOrderIsUploadedByAnotherUser = errors.New("order is uploaded by another user")
)
