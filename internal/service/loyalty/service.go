package loyalty

import (
	"context"
)

type Service interface {
	Register(ctx context.Context, login, password string) error
	Login(ctx context.Context, login, password string) (string, error)
	UploadOrder(ctx context.Context, number, userId uint) (bool, error)
	UpdateOrderAccrual(ctx context.Context, orderNumber uint) error
}
