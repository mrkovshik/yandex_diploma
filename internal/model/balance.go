package model

type GetBalanceResponse struct {
	Balance   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}
