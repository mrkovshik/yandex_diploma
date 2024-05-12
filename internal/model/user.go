package model

import (
	"database/sql"
	"time"
)

type User struct {
	ID        uint         `db:"id"`
	Login     string       `db:"login" validate:"required"`
	Password  string       `db:"password" validate:"required"`
	Balance   int          `db:"balance"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
}
