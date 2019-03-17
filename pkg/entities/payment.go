package entities

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// Payment struct represents payment entiry in the system
type Payment struct {
	// ID is an unique identifier for the payment
	ID uuid.UUID `json:"id"`

	// Account is an ID of account who requests payment information
	Account AccountID `json:"account"`

	Amount    decimal.Decimal  `json:"amount"`
	Direction PaymentDirection `json:"direction"`

	// ToAccount is a destination account ID for Outgoing payments
	ToAccount *AccountID `json:"to_account,omitempty"`

	// FromAccount is a source account ID for Incoming payments
	FromAccount *AccountID `json:"from_account,omitempty"`
}
