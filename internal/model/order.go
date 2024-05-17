package model

import (
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
	ID          uint       `db:"id" json:"-"`
	OrderNumber string     `db:"order_number" json:"number"`
	UserID      uint       `db:"user_id" json:"-"`
	Status      OrderState `db:"status" json:"status"`
	UploadedAt  time.Time  `db:"uploaded_at" json:"uploaded_at"`
	Accrual     int        `db:"accrual" json:"accrual,omitempty"`
}
