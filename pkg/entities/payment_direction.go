package entities

//go:generate stringer -type PaymentDirection

// PaymentDirection is an enum describing possible payment directions
type PaymentDirection int

const (
	// Incoming payment type is an account income
	Incoming PaymentDirection = iota

	// Outgoing payment type is an account outcome
	Outgoing
)
