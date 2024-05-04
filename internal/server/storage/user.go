package server

import "context"

type UserStorage interface {
	AddUser(ctx context.Context) error
}
