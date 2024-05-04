package postgres

import (
	"context"
	"database/sql"

	server "github.com/mrkovshik/yandex_diploma/internal/storage"
)

type postgresUserStorage struct {
	db *sql.DB
}

func NewPostgresUserStorage(db *sql.DB) server.UserStorage {
	return &postgresUserStorage{db: db}
}

func (s *postgresUserStorage) AddUser(ctx context.Context) error {
	return nil
}
