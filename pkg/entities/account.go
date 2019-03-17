package entities

import "github.com/shopspring/decimal"

// AccountID is an identifier of account
type AccountID string

// Account struct represents an user account in the system
type Account struct {
	ID       AccountID
	Currency string
	Balance  decimal.Decimal
}
