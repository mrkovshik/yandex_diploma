package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/mrkovshik/yandex_diploma/internal/app_errors"
	"github.com/mrkovshik/yandex_diploma/internal/model"
	server "github.com/mrkovshik/yandex_diploma/internal/storage"
)

type postgresOrderStorage struct {
	db *sqlx.DB
}

func NewPostgresOrderStorage(db *sqlx.DB) server.OrderStorage {
	return &postgresOrderStorage{db: db}
}

func (s *postgresOrderStorage) UploadOrder(ctx context.Context, userId, number uint) (bool, error) {
	order, err := s.GetOrderByNumber(ctx, number)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return false, err
		}
		if _, err1 := s.db.ExecContext(ctx, "INSERT INTO orders (order_number, user_id, status, uploaded_at) VALUES ($1, $2, $3, $4)",
			number, userId, model.OrderStateNew, time.Now().UTC()); err != nil {
			return false, err1
		}
		return false, nil
	}
	if order.UserId != userId {
		return false, app_errors.ErrOrderIsUploadedByAnotherUser
	}
	return true, nil
}

func (s *postgresOrderStorage) GetOrderByNumber(ctx context.Context, number uint) (order model.Order, err error) {
	err = s.db.GetContext(ctx, &order, "SELECT * FROM orders WHERE order_number=$1", number)
	return
}
