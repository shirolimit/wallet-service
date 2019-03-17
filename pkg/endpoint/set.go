package endpoint

import (
	endpoint "github.com/go-kit/kit/endpoint"
	"github.com/shirolimit/wallet-service/pkg/service"
)

// Set is a helper struct for endpoints
type Set struct {
	CreateAccountEndpoint endpoint.Endpoint
	GetAccountEndpoint    endpoint.Endpoint
	ListAccountsEndpoint  endpoint.Endpoint
	GetPaymentsEndpoint   endpoint.Endpoint
	MakePaymentEndpoint   endpoint.Endpoint
}

// NewEndpointSet creates new endpoint set
func NewEndpointSet(ws service.WalletService) Set {
	set := Set{
		CreateAccountEndpoint: MakeCreateAccountEndpoint(ws),
	}
	return set
}
