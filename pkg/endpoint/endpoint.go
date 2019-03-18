package endpoint

import (
	"context"
	"errors"

	"github.com/go-kit/kit/endpoint"
	"github.com/shirolimit/wallet-service/pkg/entities"
	"github.com/shirolimit/wallet-service/pkg/service"
)

// CreateAccountRequest is a request struct for CreateAccount method
type CreateAccountRequest struct {
	Account entities.Account
}

// CreateAccountResponse is a response struct for CreateAccount method
type CreateAccountResponse struct {
	Account entities.Account
	Error   error
}

// Failed is a Failer method implementation
func (r *CreateAccountResponse) Failed() error {
	return r.Error
}

// MakeCreateAccountEndpoint constructs CreateAccount endpoint
func MakeCreateAccountEndpoint(ws service.WalletService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(CreateAccountRequest)
		if !ok {
			return nil, errors.New("CreateAccount request type error")
		}
		err := ws.CreateAccount(ctx, req.Account)
		return CreateAccountResponse{Account: req.Account, Error: err}, nil
	}
}

// ListAccountsRequest is a request struct for ListAccounts method
type ListAccountsRequest struct {
}

// ListAccountsResponse is a response struct for ListAccounts method
type ListAccountsResponse struct {
	Accounts []entities.AccountID
	Error    error
}

// Failed is a Failer method implementation
func (r *ListAccountsResponse) Failed() error {
	return r.Error
}

// MakeListAccountsEndpoint constructs ListAccounts endpoint
func MakeListAccountsEndpoint(ws service.WalletService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		accounts, err := ws.ListAccounts(ctx)
		return ListAccountsResponse{Accounts: accounts, Error: err}, nil
	}
}

// GetAccountRequest is a request struct for GetAccount method
type GetAccountRequest struct {
	ID entities.AccountID
}

// GetAccountResponse is a response struct for GetAccount method
type GetAccountResponse struct {
	Account entities.Account
	Error   error
}

// Failed is a Failure method implementation
func (r *GetAccountResponse) Failed() error {
	return r.Error
}

// MakeGetAccountEndpoint constructs GetAccount endpoint
func MakeGetAccountEndpoint(ws service.WalletService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(GetAccountRequest)
		if !ok {
			return nil, errors.New("GetAccount request type error")
		}
		acc, err := ws.GetAccount(ctx, req.ID)
		return GetAccountResponse{Account: acc, Error: err}, nil
	}
}

// GetPaymentsRequest is a request struct for GetPayments method
type GetPaymentsRequest struct {
	AccountID entities.AccountID
}

// GetPaymentsResponse is a response struct for GetPayments method
type GetPaymentsResponse struct {
	Payments []entities.Payment
	Error    error
}

// Failed is a Failure method implementation
func (r *GetPaymentsResponse) Failed() error {
	return r.Error
}

// MakeGetPaymentsEndpoint constructs GetPayments endpoint
func MakeGetPaymentsEndpoint(ws service.WalletService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(GetPaymentsRequest)
		if !ok {
			return nil, errors.New("GetPayments request type error")
		}
		payments, err := ws.GetPayments(ctx, req.AccountID)
		return GetPaymentsResponse{Payments: payments, Error: err}, nil
	}
}

// MakePaymentRequest is a request struct for MakePayment method
type MakePaymentRequest struct {
	Payment entities.Payment
}

// MakePaymentResponse is a response struct for MakePayment method
type MakePaymentResponse struct {
	Payment entities.Payment
	Error   error
}

// Failed is a Failure method implementation
func (r *MakePaymentResponse) Failed() error {
	return r.Error
}

// MakeMakePaymentsEndpoint constructs MakePayment endpoint
func MakeMakePaymentsEndpoint(ws service.WalletService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(MakePaymentRequest)
		if !ok {
			return nil, errors.New("MakePayment request type error")
		}
		err := ws.MakePayment(ctx, req.Payment)
		return MakePaymentResponse{Payment: req.Payment, Error: err}, nil
	}
}
