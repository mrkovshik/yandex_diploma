package loyalty

import (
	"context"

	"github.com/mrkovshik/yandex_diploma/internal/model"
)

type Service interface {
	Register(ctx context.Context, login, password string) error
	Login(ctx context.Context, login, password string) (string, error)
	UploadOrder(ctx context.Context, number, userID uint) (bool, error)
	UpdateOrderAccrual(ctx context.Context, orderNumber uint) error
	GetUserOrders(ctx context.Context, userID uint) ([]model.Order, error)
	UpdatePendingOrders(ctx context.Context) error
	Withdraw(ctx context.Context, withdrawal model.Withdrawal) error
	LisUserWithdrawals(ctx context.Context, userID uint) ([]model.Withdrawal, error)
	GetBalance(ctx context.Context, userID uint) (model.GetBalanceResponse, error)
}
