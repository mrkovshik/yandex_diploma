package loyalty

import (
	"context"
	"database/sql"
	"errors"

	"github.com/mrkovshik/yandex_diploma/internal/model"
	"github.com/mrkovshik/yandex_diploma/internal/service/accrual"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"github.com/mrkovshik/yandex_diploma/internal/app_errors"
	"github.com/mrkovshik/yandex_diploma/internal/auth"
	"github.com/mrkovshik/yandex_diploma/internal/config"
)

type basicService struct {
	storage Storage
	cfg     *config.Config
	logger  *zap.SugaredLogger
}

func NewBasicService(storage Storage, cfg *config.Config, logger *zap.SugaredLogger) Service {
	return &basicService{
		storage: storage,
		cfg:     cfg,
		logger:  logger,
	}
}

func (s *basicService) Register(ctx context.Context, login, password string) error {
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return err
	}
	if err := s.storage.AddUser(ctx, login, hashedPassword); err != nil {
		return err
	}
	return nil
}

func (s *basicService) Login(ctx context.Context, login, password string) (string, error) {
	user, err := s.storage.GetUserByLogin(ctx, login)
	if err != nil {
		return "", err
	}
	if !checkPasswordHash(password, user.Password) {
		return "", app_errors.ErrInvalidPassword
	}
	authSrv := auth.NewAuthService(s.cfg.SecretKey, s.cfg.TokenExp)
	token, err := authSrv.GenerateToken(user.ID)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *basicService) UploadOrder(ctx context.Context, orderNumber, userId uint) (bool, error) {
	order, err := s.storage.GetOrderByNumber(ctx, orderNumber)

	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return false, err
		}
		if err := s.storage.UploadOrder(ctx, userId, orderNumber); err != nil {
			return false, err
		}
		return false, nil
	}

	if order.UserId != userId {
		return false, app_errors.ErrOrderIsUploadedByAnotherUser
	}
	return true, nil
}

func (s *basicService) UpdateOrderAccrual(ctx context.Context, orderNumber uint) error {
	countingSrv := accrual.NewAccrualService(s.cfg.AccrualSystemAddress)
	res, err := countingSrv.GetOrderAccrual(orderNumber)
	if err != nil {
		return err
	}
	switch res.Status {
	case model.AccrualStateInvalid:
		if err := s.storage.SetOrderStatus(ctx, orderNumber, model.OrderStateInvalid); err != nil {
			return err
		}
	case model.AccrualStateProcessing, model.AccrualStateRegistered:
		if err := s.storage.SetOrderStatus(ctx, orderNumber, model.OrderStateProcessing); err != nil {
			return err
		}
	case model.AccrualStateProcessed:
		if err := s.storage.FinalizeOrderAndUpdateBalance(ctx, orderNumber, res.Accrual); err != nil {
			return err
		}
	default:
		return errors.New("invalid accrual state")
	}

	return nil
}

func (s *basicService) GetUserOrders(ctx context.Context, userId uint) ([]model.Order, error) {
	orders, err := s.storage.GetOrdersByUserID(ctx, userId)
	if err != nil {
		return []model.Order{}, err
	}
	return orders, nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
