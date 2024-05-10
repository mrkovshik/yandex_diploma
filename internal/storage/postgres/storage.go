package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/mrkovshik/yandex_diploma/internal/app_errors"
	"github.com/mrkovshik/yandex_diploma/internal/model"
	server "github.com/mrkovshik/yandex_diploma/internal/service/loyalty"
)

type Storage struct {
	db *sqlx.DB
}

func NewStorage(db *sqlx.DB) server.Storage {
	return &Storage{db: db}
}

func (s *Storage) UploadOrder(ctx context.Context, userId, number uint) error {
	if _, err := s.db.ExecContext(ctx, "INSERT INTO orders (order_number, user_id, status, uploaded_at) VALUES ($1, $2, $3, $4)",
		number, userId, model.OrderStateNew, time.Now().UTC()); err != nil {
		return err
	}
	return nil
}

func (s *Storage) GetOrderByNumber(ctx context.Context, number uint) (order model.Order, err error) {
	err = s.db.GetContext(ctx, &order, "SELECT * FROM orders WHERE order_number=$1", number)
	return
}

func (s *Storage) SetOrderAccrual(ctx context.Context, orderNumber uint, amount int, tx *sqlx.Tx) error {
	if tx != nil {
		if _, err := tx.ExecContext(ctx, "UPDATE orders SET accrual = $1,  status = $2 WHERE order_number = $3;", amount, model.OrderStateProcessed, orderNumber); err != nil {
			return err
		}
		return nil
	}
	if _, err := s.db.ExecContext(ctx, "UPDATE orders SET accrual = $1,  status = $2 WHERE order_number = $3;", amount, model.OrderStateProcessed, orderNumber); err != nil {
		return err
	}
	return nil
}
func (s *Storage) SetOrderStatus(ctx context.Context, orderNumber uint, status model.OrderState, tx *sqlx.Tx) error {

	if tx != nil {
		if _, err := tx.ExecContext(ctx, "UPDATE orders SET status = $1 WHERE order_number = $2;", status, orderNumber); err != nil {
			return err
		}
		return nil
	}
	if _, err := s.db.ExecContext(ctx, "UPDATE orders SET status = $1 WHERE order_number = $2;", status, orderNumber); err != nil {
		return err
	}
	return nil
}

func (s *Storage) FinalizeOrderAndUpdateBalance(ctx context.Context, orderNumber uint, amount int, status model.OrderState) error {
	tx, err := s.db.Beginx()
	defer tx.Rollback() //nolint:all
	if err != nil {
		return err
	}
	if err := s.SetOrderStatus(ctx, orderNumber, status, tx); err != nil {
		return err
	}
	if err := s.SetOrderAccrual(ctx, orderNumber, amount, tx); err != nil {
		return err
	}

	if err := s.UpdateUserBalance(ctx, orderNumber, amount, tx); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (s *Storage) AddUser(ctx context.Context, login, password string) (err error) {
	_, err = s.GetUserByLogin(ctx, login)
	if err == nil {
		return app_errors.ErrUserAlreadyExists
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	if _, err = s.db.ExecContext(ctx, "INSERT INTO users (login, password, created_at) VALUES ($1, $2, $3)", login, password, time.Now().UTC()); err != nil {
		return
	}
	return
}

func (s *Storage) GetUserByLogin(ctx context.Context, login string) (user model.User, err error) {
	err = s.db.GetContext(ctx, &user, "SELECT * FROM users WHERE login=$1", login)
	return
}

func (s *Storage) GetUserByID(ctx context.Context, id uint, tx *sqlx.Tx) (user model.User, err error) {
	if tx != nil {
		err = tx.GetContext(ctx, &user, "SELECT * FROM users WHERE id=$1", id)
		return
	}
	err = s.db.GetContext(ctx, &user, "SELECT * FROM users WHERE id=$1", id)
	return
}
func (s *Storage) UpdateUserBalance(ctx context.Context, userId uint, amount int, tx *sqlx.Tx) error {

	user, err := s.GetUserByID(ctx, userId, tx)
	if err != nil {
		return err
	}
	newBalance := user.Balance + amount

	if _, err := tx.ExecContext(ctx, "UPDATE users SET balance = $1 WHERE user_id = $2;", newBalance, userId); err != nil {
		return err
	}
	return nil
}
