package loyalty

import (
	"context"
	"database/sql"
	"errors"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"github.com/mrkovshik/yandex_diploma/internal/app_errors"
	"github.com/mrkovshik/yandex_diploma/internal/auth"
	"github.com/mrkovshik/yandex_diploma/internal/config"
)

type basicService struct {
	userStorage  UserStorage
	orderStorage OrderStorage
	cfg          *config.Config
	logger       *zap.SugaredLogger
}

func NewBasicService(userStorage UserStorage, orderStorage OrderStorage, cfg *config.Config, logger *zap.SugaredLogger) Service {
	return &basicService{
		userStorage:  userStorage,
		orderStorage: orderStorage,
		cfg:          cfg,
		logger:       logger,
	}
}

func (s *basicService) Register(ctx context.Context, login, password string) error {
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return err
	}
	if err := s.userStorage.AddUser(ctx, login, hashedPassword); err != nil {
		return err
	}
	return nil
}

func (s *basicService) Login(ctx context.Context, login, password string) (string, error) {
	user, err := s.userStorage.GetUserByLogin(ctx, login)
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

func (s *basicService) UploadOrder(ctx context.Context, number, userId uint) (bool, error) {
	order, err := s.orderStorage.GetOrderByNumber(ctx, number)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return false, err
		}
		if err1 := s.orderStorage.UploadOrder(ctx, userId, number); err != nil {
			return false, err1
		}
		return false, nil
	}
	if order.UserId != userId {
		return false, app_errors.ErrOrderIsUploadedByAnotherUser
	}
	return true, nil

}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
