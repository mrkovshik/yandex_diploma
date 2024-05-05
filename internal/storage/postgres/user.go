package postgres

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	server "github.com/mrkovshik/yandex_diploma/internal/storage"
)

type postgresUserStorage struct {
	db *sqlx.DB
}

func NewPostgresUserStorage(db *sqlx.DB) server.UserStorage {
	return &postgresUserStorage{db: db}
}

func (s *postgresUserStorage) AddUser(ctx context.Context, login, password string) error {

	if _, err := s.db.ExecContext(ctx, "INSERT INTO users (login, password, created_at) VALUES ($1, $2, $3)", login, password, time.Now().UTC()); err != nil {
		return err
	}
	return nil
}
