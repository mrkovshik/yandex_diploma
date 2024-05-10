package postgres

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/mrkovshik/yandex_diploma/internal/model"
	server "github.com/mrkovshik/yandex_diploma/internal/service/loyalty"
)

type postgresOrderStorage struct {
	db *sqlx.DB
}

func NewOrderStorage(db *sqlx.DB) server.OrderStorage {
	return &postgresOrderStorage{db: db}
}

func (s *postgresOrderStorage) UploadOrder(ctx context.Context, userId, number uint) error {
	if _, err := s.db.ExecContext(ctx, "INSERT INTO orders (order_number, user_id, status, uploaded_at) VALUES ($1, $2, $3, $4)",
		number, userId, model.OrderStateNew, time.Now().UTC()); err != nil {
		return err
	}
	return nil
}

func (s *postgresOrderStorage) GetOrderByNumber(ctx context.Context, number uint) (order model.Order, err error) {
	err = s.db.GetContext(ctx, &order, "SELECT * FROM orders WHERE order_number=$1", number)
	return
}
