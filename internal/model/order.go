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
	OrderNumber uint         `db:"order_number"`
	UserId      uint         `db:"user_id"`
	Status      OrderState   `db:"status"`
	UploadedAt  time.Time    `db:"uploaded_at"`
	UpdatedAt   sql.NullTime `db:"updated_at"`
	Accrual     int          `db:"accrual"`
}
