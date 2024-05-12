package loyalty

import (
	"context"

	"github.com/mrkovshik/yandex_diploma/internal/model"
)

type Service interface {
	Register(ctx context.Context, login, password string) error
	Login(ctx context.Context, login, password string) (string, error)
	UploadOrder(ctx context.Context, number, userId uint) (bool, error)
	UpdateOrderAccrual(ctx context.Context, orderNumber uint) error
	GetUserOrders(ctx context.Context, userId uint) ([]model.Order, error)
	UpdatePendingOrders(ctx context.Context) error
}
