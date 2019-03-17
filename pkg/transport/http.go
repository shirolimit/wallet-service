package transport

import (
	"context"
	"encoding/json"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
	mux "github.com/gorilla/mux"
	"github.com/shirolimit/wallet-service/pkg/endpoint"
	"github.com/shirolimit/wallet-service/pkg/entities"
)

// NewHTTPHandler creates new HTTP handler
func NewHTTPHandler(endpoints endpoint.Set, options []httptransport.ServerOption) http.Handler {
	m := mux.NewRouter()
	makeCreateAccountHandler(m, endpoints, options)
	return m
}

func makeCreateAccountHandler(m *mux.Router, endpoints endpoint.Set, options []httptransport.ServerOption) {
	m.Methods("POST").Path("/accounts").Handler(
		httptransport.NewServer(
			endpoints.CreateAccountEndpoint,
			decodeCreateAccountRequest,
			encodeCreateAccountResponse,
			options...,
		),
	)
}

func decodeCreateAccountRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	req := endpoint.CreateAccountRequest{}
	err := json.NewDecoder(r.Body).Decode(&req.Account)
	return req, err
}

func encodeCreateAccountResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	resp, ok := response.(endpoint.CreateAccountResponse)
	if !ok || resp.Failed() != nil {
		err := resp.Failed()
		w.WriteHeader(statusCodeFromError(err))
		writeError(ctx, w, err)
		return nil
	}

	w.WriteHeader(http.StatusCreated)
	return json.NewEncoder(w).Encode(resp.Account)
}

func statusCodeFromError(err error) int {
	switch err {
	case entities.ErrAccountAlreadyExists:
		return http.StatusConflict

	case entities.ErrAccountNotFound:
		return http.StatusNotFound

	case entities.ErrDifferentCurrencies:
		return http.StatusForbidden

	case entities.ErrInsufficientFunds:
		return http.StatusPaymentRequired

	case entities.ErrPaymentAlreadyDone:
		return http.StatusConflict

	case entities.ErrRecipientNotFound:
		return http.StatusNotFound

	default:
		return http.StatusInternalServerError
	}
}

type errorWrapper struct {
	Error string `json:"error"`
}

func writeError(ctx context.Context, w http.ResponseWriter, err error) {
	json.NewEncoder(w).Encode(errorWrapper{Error: err.Error()})
}
