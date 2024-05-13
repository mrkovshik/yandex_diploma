package model

import "time"

type Withdrawal struct {
	ID          uint      `db:"id" json:"-"`
	Amount      int       `db:"amount" json:"sum"`
	ProcessedAt time.Time `db:"processed_at" json:"processed_at,omitempty"`
	OrderNumber uint      `db:"order_number" json:"order"`
	UserId      uint      `db:"user_id" json:"-"`
}

type WithdrawRequest struct {
	Sum   int    `json:"sum"`
	Order string `json:"order"`
}
