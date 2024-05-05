package model

import (
	"database/sql"
	"time"
)

type OrderState string

const (
	OrderStateNew        = OrderState("NEW")
	OrderStateProcessing = OrderState("PROCESSING")
	OrderStateInvalid    = OrderState("INVALID")
	OrderStateProcessed  = OrderState("PROCESSED")
)

type Order struct {
	ID          uint         `db:"id"`
	OrderNumber string       `db:"order_number"`
	UserId      uint         `db:"user_id"`
	Status      string       `db:"status"`
	UploadedAt  time.Time    `db:"uploaded_at"`
	UpdatedAt   sql.NullTime `db:"updated_at"`
}
