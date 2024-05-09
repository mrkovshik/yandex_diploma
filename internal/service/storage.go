package service

import (
	"context"

	"github.com/mrkovshik/yandex_diploma/internal/model"
)

type UserStorage interface {
	AddUser(ctx context.Context, name, password string) error
	GetUserByLogin(ctx context.Context, login string) (user model.User, err error)
	GetUserByID(ctx context.Context, id uint) (user model.User, err error)
}

type OrderStorage interface {
	UploadOrder(ctx context.Context, userId, number uint) error
	GetOrderByNumber(ctx context.Context, number uint) (order model.Order, err error)
}
