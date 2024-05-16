package model

type GetBalanceResponse struct {
	Balance   float64 `json:"current"`
	Withdrawn int     `json:"withdrawn"`
}
