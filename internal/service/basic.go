package service

import (
	"database/sql"
	"fmt"

	"github.com/mrkovshik/yandex_diploma/internal/config"
	"go.uber.org/zap"
)

type basicService struct {
	db     *sql.DB
	cfg    *config.Config
	logger *zap.SugaredLogger
}

func NewBasicService(db *sql.DB, cfg *config.Config, logger *zap.SugaredLogger) Service {
	return &basicService{
		db:     db,
		cfg:    cfg,
		logger: logger,
	}
}

func (s *basicService) AddUser(login, password string) error {
	fmt.Println(login, password)
	return nil
}
