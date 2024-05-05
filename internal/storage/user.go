package storage

import (
	"context"

	"github.com/mrkovshik/yandex_diploma/internal/model"
)

type UserStorage interface {
	AddUser(ctx context.Context, name, password string) error
	GetUserByLogin(ctx context.Context, login string) (user model.User, err error)
}
