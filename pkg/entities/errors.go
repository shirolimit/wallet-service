package entities

import "errors"

var (
	ErrAccountAlreadyExists       = errors.New("Account already exists")
	ErrInsufficientFunds          = errors.New("Insufficient funds to make a payment")
	ErrAccountNotFound            = errors.New("Account not found")
	ErrRecipientNotFound          = errors.New("Recipient account not found")
	ErrDifferentCurrencies        = errors.New("Payments with currency exchange are not supported")
	ErrPaymentAlreadyDone         = errors.New("Specified payment has already been completed")
	ErrDatabaseConnection         = errors.New("Database connection error")
	ErrIncomingPaymentsNotAllowed = errors.New("Incoming payments are not allowed")
	ErrWrongPaymentAmount         = errors.New("Wrong payment amount")
	ErrEmptyAccountID             = errors.New("Account ID cannot be empty")
	ErrEmptyAccountCurrency       = errors.New("Account currency cannot be empty")
	ErrNegativeBalance            = errors.New("Account balance cannot be negative")
	ErrEmptyPaymentSource         = errors.New("Payment source account cannot be empty")
	ErrEmptyPaymentDestination    = errors.New("Payment destination account cannot be empty")
	ErrEmptyPaymentID             = errors.New("Payment ID cannot be empty, use a unique GUID here")
	ErrPaymentSourceNotFound      = errors.New("Payment source account does not exist")
	ErrPaymentDestinationNotFound = errors.New("Payment destination account does not exist")
)
