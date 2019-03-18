# Wallet Service


## Summary
Is a simple service created as a test task for golang developer position.

It allows to create user accounts and carry out payments transfering money between accounts.

### Task Requirements
Use `go-kit` to build a service.

Use specified entities:
    
- Account:
        
        id: bob123 
        balance: 100.00 
        currency: USD

- Payments:

        account: bob123 
        amount: 100.00 
        to_account: alice456 
        direction: outgoing 
 
        account: alice456 
        amount: 100.00 
        from_account: bob123 
        direction: incoming 

### User stories

- Send a payment from one account to another
- See all payments
- List accounts

Added stories:
- Create new account

## Limitations
Payments between accounts in different currencies not supported.

Also there are [lots of things to do](#todo), but they were not done because of lack of time.

## Documentation

### Building

Just an ordinary go building process:

    go get ./...
    go install github.com/shirolimit/wallet-service/cmd/wallet_service

### Usage
API documentation can be found in [api.md](/docs/api.md) and [openapi.yaml](/api/openapi.yaml).

Use the following command to run service after building:

    wallet_service --connection-string=<postgres_connection_string> --http-address=":8080"

### Docker

Go to the project dir and build container:

    docker build -t wallet-service .

Run container (specify correct connection string):

    docker run -d -p 8080:8080 --env CONNECTION_STRING="user=wallet dbname=wallet_service host=127.0.0.1 password=123456 sslmode=disable" wallet-service

## TODO
Add some instrumentation:
 - tracing (i.e. Zipkin)
 - metrics (i.e. Prometheus)
 - throttling
 - circuit breaker
 - etc.

Add compose script.

Improve test coverage (endpoints, transports and middlewares are not covered at all).

Widen entities: add payment timestamps, allow currency exchange