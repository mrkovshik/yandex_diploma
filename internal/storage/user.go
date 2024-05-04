package storage

import "context"

type UserStorage interface {
	AddUser(ctx context.Context) error
}
