package transport

import (
	"context"
	"encoding/json"
	"errors"
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
	makeListAccountsHandler(m, endpoints, options)
	makeGetAccountHandler(m, endpoints, options)
	makeGetPaymentsHandler(m, endpoints, options)
	makeMakePaymentHandler(m, endpoints, options)
	return m
}

// makeCreateAccountHandler creates HTTP handler for CreateAccount requests
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

// makeListAccountsHandler creates HTTP handler for ListAccounts endpoint
func makeListAccountsHandler(m *mux.Router, endpoints endpoint.Set, options []httptransport.ServerOption) {
	m.Methods("GET").Path("/accounts").Handler(
		httptransport.NewServer(
			endpoints.ListAccountsEndpoint,
			decodeListAccountsRequest,
			encodeListAccountsResponse,
			options...,
		),
	)
}

func decodeListAccountsRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	return endpoint.ListAccountsRequest{}, nil
}

func encodeListAccountsResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	resp, ok := response.(endpoint.ListAccountsResponse)
	if !ok || resp.Failed() != nil {
		err := resp.Failed()
		w.WriteHeader(statusCodeFromError(err))
		writeError(ctx, w, err)
		return nil
	}

	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(resp.Accounts)
}

// makeGetAccountHandler creates HTTP handler for GetAccount endpoint
func makeGetAccountHandler(m *mux.Router, endpoints endpoint.Set, options []httptransport.ServerOption) {
	m.Methods("GET").Path("/accounts/{id}").Handler(
		httptransport.NewServer(
			endpoints.GetAccountEndpoint,
			decodeGetAccountRequest,
			encodeGetAccountResponse,
			options...,
		),
	)
}

func decodeGetAccountRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	req := endpoint.GetAccountRequest{
		ID: entities.AccountID(vars["id"]),
	}
	return req, nil
}

func encodeGetAccountResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	resp, ok := response.(endpoint.GetAccountResponse)
	if !ok || resp.Failed() != nil {
		err := resp.Failed()
		w.WriteHeader(statusCodeFromError(err))
		writeError(ctx, w, err)
		return nil
	}

	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(resp.Account)
}

// makeGetPaymentsHandler creates HTTP handler for GetPayments endpoint
func makeGetPaymentsHandler(m *mux.Router, endpoints endpoint.Set, options []httptransport.ServerOption) {
	m.Methods("GET").Path("/accounts/{id}/payments").Handler(
		httptransport.NewServer(
			endpoints.GetPaymentsEndpoint,
			decodeGetPaymentsRequest,
			encodeGetPaymentsResponse,
			options...,
		),
	)
}

func decodeGetPaymentsRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	req := endpoint.GetPaymentsRequest{
		AccountID: entities.AccountID(vars["id"]),
	}
	return req, nil
}

func encodeGetPaymentsResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	resp, ok := response.(endpoint.GetPaymentsResponse)
	if !ok || resp.Failed() != nil {
		err := resp.Failed()
		w.WriteHeader(statusCodeFromError(err))
		writeError(ctx, w, err)
		return nil
	}

	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(resp.Payments)
}

// makeMakePaymentHandler creates HTTP handler for MakePayment endpoint
func makeMakePaymentHandler(m *mux.Router, endpoints endpoint.Set, options []httptransport.ServerOption) {
	m.Methods("POST").Path("/accounts/{id}/payments").Handler(
		httptransport.NewServer(
			endpoints.MakePaymentEndpoint,
			decodeMakePaymentRequest,
			encodeMakePaymentResponse,
			options...,
		),
	)
}

func decodeMakePaymentRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	pathAccountID := entities.AccountID(mux.Vars(r)["id"])
	req := endpoint.MakePaymentRequest{}
	err := json.NewDecoder(r.Body).Decode(&req.Payment)
	if err != nil {
		return req, errors.New("Bad request")
	}
	if len(req.Payment.Account) == 0 {
		req.Payment.Account = pathAccountID
	}
	if req.Payment.Account != pathAccountID {
		return req, errors.New("Payment source account mismatch")
	}
	return req, nil
}

func encodeMakePaymentResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	resp, ok := response.(endpoint.MakePaymentResponse)
	if !ok || resp.Failed() != nil {
		err := resp.Failed()
		w.WriteHeader(statusCodeFromError(err))
		writeError(ctx, w, err)
		return nil
	}

	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(resp.Payment)
}

// statusCodeFromError translates error into HTTP status code
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

	case entities.ErrDatabaseConnection:
		return http.StatusInternalServerError

	case entities.ErrIncomingPaymentsNotAllowed:
		return http.StatusBadRequest

	case entities.ErrWrongPaymentAmount:
		return http.StatusBadRequest

	case entities.ErrEmptyAccountID:
		return http.StatusBadRequest

	case entities.ErrEmptyAccountCurrency:
		return http.StatusBadRequest

	case entities.ErrNegativeBalance:
		return http.StatusBadRequest

	case entities.ErrEmptyPaymentSource:
		return http.StatusBadRequest

	case entities.ErrEmptyPaymentDestination:
		return http.StatusBadRequest

	case entities.ErrEmptyPaymentID:
		return http.StatusBadRequest

	case entities.ErrPaymentSourceNotFound:
		return http.StatusNotFound

	case entities.ErrPaymentDestinationNotFound:
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
