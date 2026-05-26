# Project: supply_chain_aggregator

## Overview
Go monorepo microservices B2B supply chain aggregator. Connects SMEs with corporations via category + geospatial search.

## Tech Stack
- Language: Go 1.23
- HTTP: Echo v4 (github.com/labstack/echo/v4 v4.13.3)
- Architecture: Clean (entity / repository / usecase / delivery/http)
- DB: PostgreSQL + PostGIS (planned), currently in-memory
- Proto: protobuf for gRPC contracts

## Service Ports
- nearby-service: 8083
- transaction-service: 8084

## Service Structure (per service)
```
cmd/<service-name>/main.go
internal/
  config/config.go
  entity/
  repository/
  usecase/
  delivery/http/handler.go + router.go
go.mod  (module: supply-chain-aggregator/services/<name>)
go.sum  (copy from nearby-service if same deps)
```

## Steps Completed
- Step 1-6: monorepo, api-gateway, auth, sme, category filter, nearby-service
- Step 7: transaction-service (complete, builds, all endpoints tested)
- Step 8: payment-service (complete, builds, all endpoints tested)
- Step 9: communication-service (complete, builds, tested)

## Transaction Service
- Module: `supply-chain-aggregator/services/transaction-service`
- Endpoints:
  - POST /api/v1/transactions
  - GET  /api/v1/transactions?user_id=<optional>
  - GET  /api/v1/transactions/:id
  - PATCH /api/v1/transactions/:id/status
- Statuses: pending, paid, failed, cancelled
- ID: crypto/rand hex UUID-format
- Invoice: INV-YYYYMMDD-<seq>
- Proto: proto/transaction/transaction.proto

## Payment Service
- Module: `supply-chain-aggregator/services/payment-service`
- Port: 8085
- Deps: echo + `github.com/rabbitmq/amqp091-go v1.10.0` (go mod tidy required)
- Endpoints:
  - POST /api/v1/gateway/create-va  (invoice_number, amount, user_phone)
  - POST /api/v1/webhooks/xendit    (x-callback-token header)
- Xendit: real HTTP client; falls back to mock when XENDIT_SECRET_KEY is empty
- RabbitMQ: publisher falls back to log when RABBITMQ_URL is empty
- Webhook flow: validate token → update status → publish payment.paid event (fire-and-forget goroutine)
- RabbitMQ message shape: `{"invoice":"INV-xxx","amount":1000000,"phone":"628123456789"}`
- Env vars: XENDIT_SECRET_KEY, XENDIT_CALLBACK_TOKEN, RABBITMQ_URL

## Proto Files
- proto/nearby/nearby.proto
- proto/transaction/transaction.proto
- proto/payment/payment.proto

## Patterns
- In-memory repos with seed data (no DB yet)
- Context-aware repo methods with ctx.Done() checks
- Error wrapping: repo errors re-mapped in usecase layer
- go.sum copied from nearby-service when deps are identical
