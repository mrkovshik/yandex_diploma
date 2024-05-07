package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/mrkovshik/yandex_diploma/internal/app_errors"
	"github.com/mrkovshik/yandex_diploma/internal/auth"
	"github.com/mrkovshik/yandex_diploma/internal/config"
	"github.com/mrkovshik/yandex_diploma/internal/storage/postgres"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type basicService struct {
	db     *sqlx.DB
	cfg    *config.Config
	logger *zap.SugaredLogger
}

func NewBasicService(db *sqlx.DB, cfg *config.Config, logger *zap.SugaredLogger) Service {
	return &basicService{
		db:     db,
		cfg:    cfg,
		logger: logger,
	}
}

func (s *basicService) Register(ctx context.Context, login, password string) error {
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return err
	}
	userStorage := postgres.NewPostgresUserStorage(s.db)
	if err := userStorage.AddUser(ctx, login, hashedPassword); err != nil {
		return err
	}
	return nil
}

func (s *basicService) Login(ctx context.Context, login, password string) (string, error) {
	userStorage := postgres.NewPostgresUserStorage(s.db)
	user, err := userStorage.GetUserByLogin(ctx, login)
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
	orderStorage := postgres.NewPostgresOrderStorage(s.db)
	order, err := orderStorage.GetOrderByNumber(ctx, number)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return false, err
		}
		if err1 := orderStorage.UploadOrder(ctx, userId, number); err != nil {
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
