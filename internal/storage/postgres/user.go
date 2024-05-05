package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/mrkovshik/yandex_diploma/internal/app_errors"
	"github.com/mrkovshik/yandex_diploma/internal/model"
	server "github.com/mrkovshik/yandex_diploma/internal/storage"
)

type postgresUserStorage struct {
	db *sqlx.DB
}

func NewPostgresUserStorage(db *sqlx.DB) server.UserStorage {
	return &postgresUserStorage{db: db}
}

func (s *postgresUserStorage) AddUser(ctx context.Context, login, password string) (err error) {
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

func (s *postgresUserStorage) GetUserByLogin(ctx context.Context, login string) (user model.User, err error) {
	err = s.db.GetContext(ctx, &user, "SELECT * FROM users WHERE login=$1", login)
	return
}
