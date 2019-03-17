package entities

import "errors"

var (
	ErrAccountAlreadyExists = errors.New("Account already exists")
	ErrInsufficientFunds    = errors.New("Insufficient funds to make a payment")
	ErrAccountNotFound      = errors.New("Account not found")
	ErrRecipientNotFound    = errors.New("Recipient account not found")
	ErrDifferentCurrencies  = errors.New("Payments with currency exchange are not supported")
	ErrPaymentAlreadyDone   = errors.New("Specified payment has already been completed")
)
