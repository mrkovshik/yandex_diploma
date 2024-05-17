package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	server "github.com/mrkovshik/yandex_diploma/api"

	"github.com/jmoiron/sqlx"
	"github.com/mrkovshik/yandex_diploma/internal/apperrors"
	"github.com/mrkovshik/yandex_diploma/internal/model"
)

type Storage struct {
	db *sqlx.DB
}

func NewStorage(db *sqlx.DB) server.Storage {
	return &Storage{db: db}
}

func (s *Storage) UploadOrder(ctx context.Context, userID, number uint) error {
	if _, err := s.db.ExecContext(ctx, "INSERT INTO orders (order_number, user_id, status, uploaded_at) VALUES ($1, $2, $3, $4)",
		number, userID, model.OrderStateNew, time.Now().UTC()); err != nil {
		return err
	}
	return nil
}

func (s *Storage) SetOrderStatus(ctx context.Context, orderNumber uint, status model.OrderState) error {
	if _, err := s.db.ExecContext(ctx, "UPDATE orders SET status = $1 WHERE order_number = $2;", status, orderNumber); err != nil {
		return err
	}
	return nil
}

func (s *Storage) GetOrderByNumber(ctx context.Context, number uint) (order model.Order, err error) {
	err = s.db.GetContext(ctx, &order, "SELECT * FROM orders WHERE order_number=$1", number)
	return
}

func (s *Storage) FinalizeOrderAndUpdateBalance(ctx context.Context, orderNumber uint, amount int) error {
	tx, err := s.db.Beginx()
	defer tx.Rollback() //nolint:all
	if err != nil {
		return err
	}
	if err := s.setOrderStatusTx(ctx, orderNumber, model.OrderStateProcessed, tx); err != nil {
		return err
	}
	if err := s.setOrderAccrualTx(ctx, orderNumber, amount, tx); err != nil {
		return err
	}
	if err := s.updateUserBalanceByOrderNumberTx(ctx, orderNumber, amount, tx); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (s *Storage) AddUser(ctx context.Context, login, password string) (uint, error) {
	_, err := s.GetUserByLogin(ctx, login)
	if err == nil {
		return 0, apperrors.ErrUserAlreadyExists
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return 0, err
	}
	if _, err := s.db.ExecContext(ctx, "INSERT INTO users (login, password, created_at) VALUES ($1, $2, $3)", login, password, time.Now().UTC()); err != nil {
		return 0, err
	}
	user, err1 := s.GetUserByLogin(ctx, login)
	if err1 != nil {
		return 0, err1
	}
	return user.ID, nil
}

func (s *Storage) GetUserByLogin(ctx context.Context, login string) (user model.User, err error) {
	err = s.db.GetContext(ctx, &user, "SELECT * FROM users WHERE login=$1", login)
	return
}

func (s *Storage) GetUserByID(ctx context.Context, id uint) (user model.User, err error) {
	err = s.db.GetContext(ctx, &user, "SELECT * FROM users WHERE id=$1", id)
	return
}

func (s *Storage) GetOrdersByUserID(ctx context.Context, userID uint) (orders []model.Order, err error) {
	err = s.db.SelectContext(ctx, &orders, "SELECT * FROM orders WHERE user_id=$1", userID)
	return
}

func (s *Storage) GetPendingOrders(ctx context.Context) (orders []uint, err error) {
	err = s.db.SelectContext(ctx, &orders, "SELECT order_number FROM orders WHERE status IN ($1, $2)", model.OrderStateNew, model.OrderStateProcessing)
	return
}

func (s *Storage) ProcessWithdrawal(ctx context.Context, withdrawal model.Withdrawal) error {
	tx, err := s.db.Beginx()
	defer tx.Rollback() //nolint:all
	if err != nil {
		return err
	}
	if err := s.addWithdrawalTx(ctx, withdrawal, tx); err != nil {
		return err
	}
	if err := s.updateUserBalanceByUserIDTx(ctx, withdrawal.UserID, -withdrawal.Amount, tx); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (s *Storage) GetWithdrawalsSumByUserID(ctx context.Context, userID uint) (int, error) {
	var sums []int
	err := s.db.SelectContext(ctx, &sums, "SELECT SUM(amount) FROM withdrawals WHERE withdrawals.user_id = $1 group by user_id", userID)
	if err != nil {
		return 0, err
	}
	if len(sums) == 0 {
		return 0, sql.ErrNoRows
	}

	return sums[0], nil
}

func (s *Storage) GetWithdrawalsByUserID(ctx context.Context, userID uint) (withdrawals []model.Withdrawal, err error) {
	err = s.db.SelectContext(ctx, &withdrawals, "SELECT * FROM withdrawals WHERE user_id = $1", userID)
	return
}

func (s *Storage) updateUserBalanceByOrderNumberTx(ctx context.Context, orderNumber uint, amount int, tx *sqlx.Tx) error {

	user, err := s.getUserByOrderNumberTx(ctx, orderNumber, tx)
	if err != nil {
		return err
	}
	newBalance := user.Balance + float64(amount)
	if newBalance < 0 {
		return apperrors.ErrNotEnoughFunds
	}

	if _, err := tx.ExecContext(ctx, "UPDATE users SET balance = $1 WHERE id = $2;", newBalance, user.ID); err != nil {
		return err
	}
	return nil
}

func (s *Storage) updateUserBalanceByUserIDTx(ctx context.Context, userID uint, amount int, tx *sqlx.Tx) error {
	user, err := s.getUserByUserIDTx(ctx, userID, tx)
	if err != nil {
		return err
	}

	newBalance := user.Balance + float64(amount)
	if newBalance < 0 {
		return apperrors.ErrNotEnoughFunds
	}

	if _, err := tx.ExecContext(ctx, "UPDATE users SET balance = $1 WHERE id = $2;", newBalance, user.ID); err != nil {
		return err
	}
	return nil
}

func (s *Storage) setOrderStatusTx(ctx context.Context, orderNumber uint, status model.OrderState, tx *sqlx.Tx) error {
	_, err := tx.ExecContext(ctx, "UPDATE orders SET status = $1 WHERE order_number = $2;", status, orderNumber)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) setOrderAccrualTx(ctx context.Context, orderNumber uint, amount int, tx *sqlx.Tx) error {
	if _, err := tx.ExecContext(ctx, "UPDATE orders SET accrual = $1,  status = $2 WHERE order_number = $3;", amount, model.OrderStateProcessed, orderNumber); err != nil {
		return err
	}
	return nil
}
func (s *Storage) getUserByOrderNumberTx(ctx context.Context, id uint, tx *sqlx.Tx) (user model.User, err error) {
	err = tx.GetContext(ctx, &user, "SELECT u.id, login, password, created_at, balance FROM users u join orders o on u.id = o.user_id WHERE o.order_number=$1", id)
	return
}

func (s *Storage) getUserByUserIDTx(ctx context.Context, userID uint, tx *sqlx.Tx) (user model.User, err error) {
	err = tx.GetContext(ctx, &user, "SELECT * FROM users  WHERE id =$1", userID)
	return
}

func (s *Storage) addWithdrawalTx(ctx context.Context, withdrawal model.Withdrawal, tx *sqlx.Tx) error {
	if _, err := tx.ExecContext(ctx, "INSERT INTO withdrawals (amount, processed_at, order_number, user_id) VALUES ($1, $2, $3, $4)", withdrawal.Amount, time.Now().UTC(), withdrawal.OrderNumber, withdrawal.UserID); err != nil {
		return err
	}
	return nil
}
