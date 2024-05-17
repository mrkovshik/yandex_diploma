package api

import (
	"context"

	"github.com/mrkovshik/yandex_diploma/internal/model"
)

type Storage interface {
	AddUser(ctx context.Context, login, password string) (uint, error)
	GetUserByLogin(ctx context.Context, login string) (user model.User, err error)
	GetUserByID(ctx context.Context, id uint) (user model.User, err error)
	UploadOrder(ctx context.Context, userID uint, orderNumber string) error
	GetOrderByNumber(ctx context.Context, orderNumber string) (order model.Order, err error)
	FinalizeOrderAndUpdateBalance(ctx context.Context, orderNumber string, amount int) error
	SetOrderStatus(ctx context.Context, orderNumber string, status model.OrderState) error
	GetOrdersByUserID(ctx context.Context, userID uint) ([]model.Order, error)
	GetPendingOrders(ctx context.Context) (orders []string, err error)
	ProcessWithdrawal(ctx context.Context, withdrawal model.Withdrawal) error
	GetWithdrawalsSumByUserID(ctx context.Context, userID uint) (sum float64, err error)
	GetWithdrawalsByUserID(ctx context.Context, userID uint) (withdrawals []model.Withdrawal, err error)
}
