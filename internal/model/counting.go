package model

type AccrualState string

const (
	AccrualStateRegistered = AccrualState("REGISTERED")
	AccrualStateProcessing = AccrualState("PROCESSING")
	AccrualStateInvalid    = AccrualState("INVALID")
	AccrualStateProcessed  = AccrualState("PROCESSED")
)
