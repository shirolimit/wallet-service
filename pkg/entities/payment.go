package entities

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// Payment struct represents payment entiry in the system
type Payment struct {
	// ID is an unique identifier for the payment
	ID uuid.UUID

	// Account is an ID of account who requests payment information
	Account AccountID

	Amount    decimal.Decimal
	Direction PaymentDirection

	// ToAccount is a destination account ID for Outgoing payments
	ToAccount *AccountID

	// FromAccount is a source account ID for Incoming payments
	FromAccount *AccountID
}
