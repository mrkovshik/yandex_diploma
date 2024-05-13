package loyalty

import (
	"context"

	"github.com/mrkovshik/yandex_diploma/internal/model"
)

type Storage interface {
	AddUser(ctx context.Context, name, password string) error
	GetUserByLogin(ctx context.Context, login string) (user model.User, err error)
	GetUserByID(ctx context.Context, id uint) (user model.User, err error)

	UploadOrder(ctx context.Context, userId, orderNumber uint) error
	GetOrderByNumber(ctx context.Context, orderNumber uint) (order model.Order, err error)
	FinalizeOrderAndUpdateBalance(ctx context.Context, orderNumber uint, amount int) error
	SetOrderStatus(ctx context.Context, orderNumber uint, status model.OrderState) error
	GetOrdersByUserID(ctx context.Context, userId uint) ([]model.Order, error)
	GetPendingOrders(ctx context.Context) (orders []uint, err error)
	ProcessWithdrawal(ctx context.Context, withdrawal model.Withdrawal) error
}
