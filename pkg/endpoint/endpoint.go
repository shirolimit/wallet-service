package endpoint

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/shirolimit/wallet-service/pkg/entities"
	"github.com/shirolimit/wallet-service/pkg/service"
)

type CreateAccountRequest struct {
	Account entities.Account
}

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
		// TODO: check for casting result
		req := request.(CreateAccountRequest)
		err := ws.CreateAccount(ctx, req.Account)
		return CreateAccountResponse{Account: req.Account, Error: err}, nil
	}
}
