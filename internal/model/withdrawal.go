package model

import "time"

type Withdrawal struct {
	ID          uint      `db:"id" json:"-"`
	Amount      int       `db:"amount" json:"sum"`
	ProcessedAt time.Time `db:"processed_at" json:"processed_at"`
	OrderNumber string    `db:"order_number" json:"order"`
	UserID      uint      `db:"user_id" json:"-"`
}
