package model

import (
	"database/sql"
	"time"
)

type User struct {
	ID        uint         `db:"id"`
	Login     string       `db:"login"`
	Password  string       `db:"password"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
}
