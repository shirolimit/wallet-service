# Wallet Service


## Summary
Is a simple service created as a test task for golang developer position.

It allows to create user accounts and carry out payments transfering money between accounts.

### Requirements
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

### Usage
API documentation can be found in [api.md](/docs/api.md) and [openapi.yaml](/api/openapi.yaml).

## TODO
Add some instrumentation:
 - tracing (i.e. Zipkin)
 - metrics (i.e. Prometheus)
 - throttling
 - circuit breaker
 - etc.

Improve test coverage (endpoints, transports and middlewares are not covered at all).

Widen entities: add payment timestamps, allow currency exchange