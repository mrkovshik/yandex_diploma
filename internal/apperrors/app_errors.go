package apperrors

import "errors"

var (
	ErrUserAlreadyExists = errors.New("user is already exist")
	ErrInvalidPassword   = errors.New("password is invalid")

	ErrOrderIsUploadedByAnotherUser = errors.New("order is uploaded by another user")

	ErrNoSuchOrder         = errors.New("order is not registered in loyalty program")
	ErrInvalidResponseCode = errors.New("response code is invalid")
	ErrTooManyRetrials     = errors.New("quota exceeded")

	ErrNotEnoughFunds = errors.New("not enough funds on user's balance")
)
