package model

type AccrualState string

const (
	AccrualStateRegistered = AccrualState("REGISTERED")
	AccrualStateProcessing = AccrualState("PROCESSING")
	AccrualStateInvalid    = AccrualState("INVALID")
	AccrualStateProcessed  = AccrualState("PROCESSED")
)

type AccrualResponse struct {
	Order   string       `json:"order" uri:"order" binding:"required"`
	Status  AccrualState `json:"status"`
	Accrual float64      `json:"accrual"`
}
