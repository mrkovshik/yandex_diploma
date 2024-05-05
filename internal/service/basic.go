package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"

	"github.com/jmoiron/sqlx"
	"github.com/mrkovshik/yandex_diploma/internal/config"
	"github.com/mrkovshik/yandex_diploma/internal/storage/postgres"
	"go.uber.org/zap"
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
	hasher := sha256.New()
	hasher.Write([]byte(password))
	hashSum := hasher.Sum(nil)
	hashedPassword := hex.EncodeToString(hashSum)
	userStorage := postgres.NewPostgresUserStorage(s.db)
	if err := userStorage.AddUser(ctx, login, hashedPassword); err != nil {
		return err
	}
	return nil
}
