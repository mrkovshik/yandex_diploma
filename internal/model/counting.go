package model

type CountingState string

const (
	CountingStateRegistered = CountingState("REGISTERED")
	CountingStateProcessing = CountingState("PROCESSING")
	CountingStateInvalid    = CountingState("INVALID")
	CountingStateProcessed  = CountingState("PROCESSED")
)
