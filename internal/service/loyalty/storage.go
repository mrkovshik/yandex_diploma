package loyalty

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/mrkovshik/yandex_diploma/internal/model"
)

type Storage interface {
	AddUser(ctx context.Context, name, password string) error
	GetUserByLogin(ctx context.Context, login string) (user model.User, err error)
	GetUserByID(ctx context.Context, id uint, tx *sqlx.Tx) (user model.User, err error)
	UpdateUserBalance(ctx context.Context, orderNumber uint, amount int, tx *sqlx.Tx) error

	UploadOrder(ctx context.Context, userId, orderNumber uint) error
	GetOrderByNumber(ctx context.Context, orderNumber uint) (order model.Order, err error)
	SetOrderAccrual(ctx context.Context, orderNumber uint, amount int, tx *sqlx.Tx) error
	SetOrderStatus(ctx context.Context, orderNumber uint, status model.OrderState, tx *sqlx.Tx) error
}
