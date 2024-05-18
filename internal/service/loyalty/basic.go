package loyalty

import (
	"context"
	"database/sql"
	"errors"

	"github.com/mrkovshik/yandex_diploma/internal/service"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"github.com/mrkovshik/yandex_diploma/api"
	"github.com/mrkovshik/yandex_diploma/internal/apperrors"
	"github.com/mrkovshik/yandex_diploma/internal/auth"
	"github.com/mrkovshik/yandex_diploma/internal/config"
	"github.com/mrkovshik/yandex_diploma/internal/model"
)

const workersQty = 2

type (
	basicService struct {
		storage service.Storage
		cfg     *config.Config
		accrual AccrualService
		Logger  *zap.SugaredLogger
	}
)

func NewBasicService(storage service.Storage, accrual AccrualService, cfg *config.Config, logger *zap.SugaredLogger) api.Service {
	return &basicService{
		storage: storage,
		accrual: accrual,
		cfg:     cfg,
		Logger:  logger,
	}
}

func (s *basicService) Register(ctx context.Context, login, password string) (string, error) {
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return "", err
	}
	userID, err := s.storage.AddUser(ctx, login, hashedPassword)
	if err != nil {
		return "", err
	}
	authSrv := auth.NewAuthService(s.cfg.SecretKey, s.cfg.TokenExp)
	token, err := authSrv.GenerateToken(userID)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *basicService) Login(ctx context.Context, login, password string) (string, error) {
	user, err := s.storage.GetUserByLogin(ctx, login)
	if err != nil {
		return "", err
	}
	if !checkPasswordHash(password, user.Password) {
		return "", apperrors.ErrInvalidPassword
	}
	authSrv := auth.NewAuthService(s.cfg.SecretKey, s.cfg.TokenExp)
	token, err := authSrv.GenerateToken(user.ID)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *basicService) UploadOrder(ctx context.Context, orderNumber string, userID uint) (bool, error) {
	order, err := s.storage.GetOrderByNumber(ctx, orderNumber)

	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return false, err
		}
		if err := s.storage.UploadOrder(ctx, userID, orderNumber); err != nil {
			return false, err
		}
		return false, nil
	}

	if order.UserID != userID {
		return false, apperrors.ErrOrderIsUploadedByAnotherUser
	}
	return true, nil
}

func (s *basicService) UpdateOrderAccrual(ctx context.Context, orderNumber string) error {
	res, err := s.accrual.GetOrderAccrual(orderNumber)
	if err != nil {
		return err
	}
	switch res.Status {
	case model.AccrualStateInvalid:
		if err := s.storage.SetOrderStatus(ctx, orderNumber, model.OrderStateInvalid); err != nil {
			return err
		}
		s.Logger.Debugf("updated order %v state = INVALID", orderNumber)
	case model.AccrualStateProcessing, model.AccrualStateRegistered:
		if err := s.storage.SetOrderStatus(ctx, orderNumber, model.OrderStateProcessing); err != nil {
			return err
		}
		s.Logger.Debugf("updated order %v state = PROCESSING", orderNumber)
	case model.AccrualStateProcessed:
		if err := s.storage.FinalizeOrderAndUpdateBalance(ctx, orderNumber, res.Accrual); err != nil {
			return err
		}
		s.Logger.Debugf("updated order %v with amount = %v and state = PROCESSED", orderNumber, res.Accrual)
	default:
		return errors.New("invalid accrual state")
	}

	return nil
}

func (s *basicService) UpdatePendingOrders(ctx context.Context) error {
	s.Logger.Debug("updating orders started")
	orders, err := s.storage.GetPendingOrders(ctx)
	if err != nil {
		return err
	}
	jobs := make(chan string, len(orders))
	for w := 1; w <= workersQty; w++ {
		go s.worker(ctx, w, jobs)
	}
	for _, id := range orders {
		jobs <- id
	}
	close(jobs)
	return nil
}

func (s *basicService) GetUserOrders(ctx context.Context, userID uint) ([]model.Order, error) {
	orders, err := s.storage.GetOrdersByUserID(ctx, userID)
	if err != nil {
		return []model.Order{}, err
	}
	return orders, nil
}

func (s *basicService) Withdraw(ctx context.Context, withdrawal model.Withdrawal) error {
	if err := s.storage.ProcessWithdrawal(ctx, withdrawal); err != nil {
		return err
	}
	return nil
}

func (s *basicService) GetBalance(ctx context.Context, userID uint) (model.GetBalanceResponse, error) {
	user, err := s.storage.GetUserByID(ctx, userID)
	if err != nil {
		return model.GetBalanceResponse{}, err
	}
	withdrawn, err1 := s.storage.GetWithdrawalsSumByUserID(ctx, userID)
	if err1 != nil {
		return model.GetBalanceResponse{}, err1
	}
	return model.GetBalanceResponse{
		Balance:   user.Balance,
		Withdrawn: withdrawn,
	}, nil
}

func (s *basicService) ListUserWithdrawals(ctx context.Context, userID uint) ([]model.Withdrawal, error) {

	withdrawals, err := s.storage.GetWithdrawalsByUserID(ctx, userID)
	if err != nil {
		return []model.Withdrawal{}, err
	}
	return withdrawals, nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (s *basicService) worker(ctx context.Context, workerID int, jobs <-chan string) {
	for orderNumber := range jobs {

		if err := s.UpdateOrderAccrual(ctx, orderNumber); err != nil {
			s.Logger.Errorf("failed to update order #%v by worker #%v", orderNumber, workerID)
			return
		}
	}
}
