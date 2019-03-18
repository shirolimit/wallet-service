package entities

import (
	"encoding/json"

	"github.com/shopspring/decimal"
)

// AccountID is an identifier of account
type AccountID string

// Account struct represents an user account in the system
type Account struct {
	ID       AccountID       `json:"id"`
	Currency string          `json:"currency"`
	Balance  decimal.Decimal `json:"balance"`
}

// String implements Stringer interface for logging
func (a Account) String() string {
	if data, err := json.Marshal(a); err == nil {
		return string(data)
	}
	return "account"
}
