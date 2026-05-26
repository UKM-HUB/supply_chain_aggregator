# supply_chain_aggregator

`supply_chain_aggregator` is a B2B supply chain aggregator platform designed to connect small and medium enterprises with large corporations. The application helps corporations discover nearby SMEs based on business needs, product/service categories, and geographic location using latitude and longitude.

## Project Summary

This project will be built as a Go-based microservices application inside a monorepo. Each service owns a clear business capability and communicates through HTTP APIs, gRPC, and asynchronous events where needed.

The main goal is to provide a platform where:

- **Large corporations** can search for SMEs based on required categories, location, and availability.
- **SMEs/UMKM** can register their business profile, location, products, and supply capability.
- **The system** can recommend the nearest matching SMEs using latitude and longitude.
- **Transactions and payments** can be tracked from order creation until payment completion.
- **Notifications and reports** can be generated from transaction events.

## Core Features

- **Authentication and Authorization**
  - JWT-based login, refresh token, logout, and protected API access.
  - User roles can separate access between admin, corporation, and SME users.
  - Auth validation can be exposed through gRPC for internal service-to-service authorization.

- **SME Registration**
  - SMEs can register using business identity, phone number, password, and location data.
  - Registration flow supports OTP verification before creating the user account.

- **Category-Based Discovery**
  - Corporations can filter SMEs based on categories that match their procurement needs.
  - Categories represent business sectors, products, or services offered by SMEs.

- **Nearby SME Search**
  - Corporations can search nearby SMEs using latitude and longitude.
  - The nearby service can use PostgreSQL with PostGIS for geographic queries.
  - Example endpoint:

```http
GET /api/v1/nearby/umkm?lat=-6.2&lng=106.8
```

- **Transaction Management**
  - The transaction service manages invoices, amounts, payment status, payment method, and transaction history.
  - Transaction endpoints include create, list, and detail APIs.

- **Payment Gateway Integration**
  - Payment flow supports virtual account creation through Xendit.
  - Xendit webhooks update transaction status after payment.
  - Successful payment events are published to RabbitMQ.

- **Communication Service**
  - Consumes RabbitMQ events such as `payment.paid`.
  - Sends WhatsApp notifications after successful payment.

- **Reporting Service**
  - Provides daily, monthly, and exportable reports.
  - Reports include total transactions, total paid amount, and pending transaction count.

## Technology Stack

- **Language:** Go
- **Architecture:** Microservices in a monorepo
- **HTTP Framework:** Echo or similar Go HTTP framework
- **Internal Communication:** gRPC
- **Database:** PostgreSQL
- **ORM:** GORM
- **Geospatial Search:** PostgreSQL PostGIS
- **Authentication:** JWT
- **API Contract:** OpenAPI/Swagger
- **Message Broker:** RabbitMQ
- **Payment Gateway:** Xendit
- **Notification Channel:** WhatsApp service integration

## Proposed Monorepo Structure

```text
supply_chain_aggregator/
├── services/
│   ├── api-gateway/
│   ├── auth-service/
│   ├── user-service/
│   ├── sme-service/
│   ├── nearby-service/
│   ├── transaction-service/
│   ├── payment-service/
│   ├── report-service/
│   └── communication-service/
├── proto/
│   ├── auth/
│   ├── user/
│   ├── sme/
│   ├── nearby/
│   └── transaction/
├── contracts/
│   └── openapi/
├── pkg/
│   ├── config/
│   ├── database/
│   ├── jwt/
│   ├── logger/
│   └── rabbitmq/
├── migrations/
├── deployments/
├── scripts/
├── docs/
└── README.md
```

## Service Overview

- **API Gateway**
  - Public entry point for `/api/v1`.
  - Handles routing, request validation, JWT middleware, and Swagger/OpenAPI exposure.
  - Calls internal services through gRPC or HTTP depending on service responsibility.

- **Auth Service**
  - Handles register, login, refresh token, logout, password hashing, and token generation.
  - Exposes gRPC token validation for other services.

- **User Service**
  - Manages user profile data and role-based user information.
  - Supports paginated list endpoints for users.

- **SME Service**
  - Manages SME/UMKM profiles, categories, business details, product/service availability, and coordinates.
  - Provides category-filtered SME search.

- **Nearby Service**
  - Handles geographic search using latitude and longitude.
  - Uses PostGIS distance queries to return nearest SMEs.

- **Transaction Service**
  - Manages order and transaction records.
  - Tracks invoice number, amount, payment method, and payment status.

- **Payment Service**
  - Integrates with Xendit for virtual account creation.
  - Handles payment webhook callbacks.
  - Publishes payment events to RabbitMQ.

- **Communication Service**
  - Consumes transaction/payment events.
  - Sends WhatsApp notifications to related users.

- **Report Service**
  - Generates daily, monthly, and exportable business reports.

## Main API Structure

```text
/api/v1/
├── auth
├── users
├── umkm
├── factories
├── transactions
├── reports
├── nearby
├── communication
└── gateway
```

## Example Business Flow

```text
Corporation logs in
   ↓
Corporation searches SMEs by category and location
   ↓
Nearby service returns closest matching SMEs
   ↓
Corporation creates transaction/order
   ↓
Payment service creates virtual account through Xendit
   ↓
Customer completes payment
   ↓
Xendit sends webhook
   ↓
Transaction status is updated
   ↓
RabbitMQ publishes payment.paid event
   ↓
Communication service sends WhatsApp notification
   ↓
Report service updates reporting data
```

## API Contract

All public HTTP APIs should be documented using OpenAPI/Swagger contracts. The contracts should define:

- **Request and response schemas**
- **Authentication requirements**
- **Error response format**
- **Pagination format**
- **Category filter query parameters**
- **Latitude and longitude search parameters**
- **Transaction and payment webhook payloads**

## Nearby Search Concept

The nearby search can be implemented with PostgreSQL PostGIS using coordinate-based distance calculation.

Example query concept:

```sql
SELECT *,
ST_Distance(
  location,
  ST_MakePoint(106.8, -6.2)
)
FROM umkms
ORDER BY location <-> ST_MakePoint(106.8, -6.2)
LIMIT 10;
```

This allows corporations to find the closest SMEs that match their category requirements.

## Step-by-Step Implementation Plan

### Step 1: Prepare the Monorepo

Create the base monorepo structure for all services, shared packages, protobuf contracts, OpenAPI contracts, migrations, deployment files, and documentation.

Expected folders:

```text
services/
proto/
contracts/openapi/
pkg/
migrations/
deployments/
scripts/
docs/
```

Each service should follow a clean Go structure:

```text
cmd/
internal/
├── delivery/
│   ├── http/
│   └── grpc/
├── usecase/
├── repository/
├── entity/
├── middleware/
├── config/
└── helper/
```

### Step 2: Build the API Gateway

The API Gateway is the public entry point for clients. It exposes `/api/v1` endpoints and forwards requests to internal services.

Responsibilities:

- **Routing:** Route public HTTP requests to the correct service.
- **JWT Middleware:** Protect private endpoints using JWT validation.
- **OpenAPI/Swagger:** Serve API documentation for public contracts.
- **Request Validation:** Validate request payloads before forwarding.
- **Service Communication:** Communicate with internal services using gRPC or HTTP.

Main API groups:

```text
/api/v1/auth
/api/v1/users
/api/v1/umkm
/api/v1/factories
/api/v1/transactions
/api/v1/reports
/api/v1/nearby
/api/v1/communication
/api/v1/gateway
```

### Step 3: Build the Auth Service

The Auth Service handles user identity and access control.

Endpoints based on the docs:

```http
POST /api/v1/auth/login
POST /api/v1/auth/register
POST /api/v1/auth/refresh
POST /api/v1/auth/logout
```

Registration flow:

```text
Register
   ↓
Send OTP
   ↓
Verify OTP
   ↓
Create User
```

Login response concept:

```json
{
  "token": "jwt-token",
  "refresh_token": "refresh-token",
  "user": {
    "id": 1,
    "name": "Budi",
    "role": "ADMIN"
  }
}
```

Internal gRPC contract:

```text
ValidateToken(token) returns is_valid, user_id, role
```

This allows other services to validate JWT tokens without duplicating authentication logic.

### Step 4: Build the SME Service

The SME Service manages SME/UMKM data and business capabilities.

Responsibilities:

- **Profile Management:** Store SME name, phone, address, latitude, longitude, and owner data.
- **Category Management:** Attach categories to SMEs based on their products or services.
- **Availability Data:** Store products, stock, production capacity, or service availability.
- **Search Support:** Provide category-filtered SME data to the Nearby Service or API Gateway.

Example SME registration request from the docs:

```json
{
  "name": "UMKM Maju",
  "phone": "628123456789",
  "password": "123456"
}
```

For the final application, the SME profile should also include:

- **Category IDs:** Used for corporation needs filtering.
- **Latitude and longitude:** Used for nearby search.
- **Business description:** Used to explain SME capabilities.
- **Status:** Used to control active or inactive SME visibility.

### Step 5: Build the Category Filter

Categories are used by corporations to find SMEs that match procurement needs.

Example category use cases:

- **Food supplier**
- **Packaging**
- **Textile**
- **Raw material**
- **Logistics**
- **Manufacturing support**

Search behavior:

```text
Corporation selects category
   ↓
System filters SMEs by matching category
   ↓
System combines category result with nearby search
   ↓
Corporation receives nearest matching SMEs
```

Suggested query parameters:

```http
GET /api/v1/nearby/umkm?lat=-6.2&lng=106.8&category_id=food&radius_km=10
```

### Step 6: Build the Nearby Service

The Nearby Service handles geospatial discovery of SMEs.

Docs endpoint:

```http
GET /api/v1/nearby/umkm?lat=-6.2&lng=106.8
```

Recommended final endpoint:

```http
GET /api/v1/nearby/umkm?lat=-6.2&lng=106.8&category_id=food&radius_km=10&limit=10
```

Database requirement:

- **PostgreSQL** stores SME data.
- **PostGIS** enables distance-based search.
- **GORM** manages database access from Go services.

Search flow:

```text
Receive lat/lng from corporation
   ↓
Validate coordinates and category filter
   ↓
Query SMEs using PostGIS distance calculation
   ↓
Sort result by nearest distance
   ↓
Return paginated SME list
```

### Step 7: Build the Transaction Service

The Transaction Service manages order and payment records.

Table concept from the docs:

```sql
CREATE TABLE transactions (
    id UUID PRIMARY KEY,
    invoice_number VARCHAR(100),
    user_id UUID,
    amount NUMERIC,
    status VARCHAR(20),
    payment_method VARCHAR(50),
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
```

Endpoints:

```http
POST /api/v1/transactions
GET  /api/v1/transactions
GET  /api/v1/transactions/:id
```

Transaction status examples:

- **pending:** Transaction created but not paid.
- **paid:** Payment has been completed.
- **failed:** Payment failed or expired.
- **cancelled:** Transaction was cancelled.

### Step 8: Build the Payment Gateway Service

The Payment Gateway Service integrates with Xendit.

Endpoints from the docs:

```http
POST /api/v1/gateway/create-va
POST /api/v1/webhooks/xendit
```

Payment flow:

```text
Client
   ↓
Backend
   ↓
Xendit API
   ↓
VA Created
   ↓
Customer Transfer
   ↓
Xendit Webhook
   ↓
Webhook Handler
   ↓
Update Transaction
   ↓
Publish RabbitMQ
```

Webhook security:

- **Callback Token:** Validate `x-callback-token` from Xendit.
- **Environment Variable:** Store `XENDIT_CALLBACK_TOKEN` securely.
- **No Hardcoded Secret:** Do not put payment secrets in source code.

### Step 9: Build the Communication Service

The Communication Service consumes payment events and sends notifications.

RabbitMQ queue from the docs:

```text
payment.paid
```

Message concept:

```json
{
  "invoice": "INV-001",
  "amount": 1000000,
  "phone": "628123456789"
}
```

Flow:

```text
Transaction PAID
   ↓
RabbitMQ Publish
   ↓
Notification Service Consume
   ↓
Send WhatsApp
```

WhatsApp message concept:

```text
Pembayaran berhasil diterima.

Invoice: INV-001
Nominal: Rp1.000.000

Terima kasih.
```

### Step 10: Build the Report Service

The Report Service provides business reporting for admins and corporations.

Endpoints:

```http
GET /api/v1/reports/daily
GET /api/v1/reports/monthly
GET /api/v1/reports/export
```

Response concept:

```json
{
  "total_transaction": 100,
  "total_paid": 50000000,
  "total_pending": 5
}
```

Report data should be generated from transaction records and payment statuses.

### Step 11: Define OpenAPI/Swagger Contracts

Each public HTTP API should have a Swagger/OpenAPI contract under:

```text
contracts/openapi/
```

Recommended contract files:

```text
auth.yaml
users.yaml
sme.yaml
nearby.yaml
transactions.yaml
payments.yaml
reports.yaml
```

Each contract should define:

- **Endpoint path and method**
- **Request body schema**
- **Response body schema**
- **JWT security scheme**
- **Query parameters**
- **Pagination parameters**
- **Error response format**

### Step 12: Define gRPC Contracts

Internal service-to-service communication should use protobuf files under:

```text
proto/
```

Recommended protobuf modules:

```text
proto/auth/auth.proto
proto/sme/sme.proto
proto/nearby/nearby.proto
proto/transaction/transaction.proto
proto/payment/payment.proto
```

The Auth Service should expose token validation through gRPC so other services can verify identity and role data.

### Step 13: Database and Migration Strategy

Use PostgreSQL as the main database and GORM as the ORM layer.

Recommended core tables:

- **users:** Stores login identity, role, password hash, and user metadata.
- **smes:** Stores SME business profile and coordinates.
- **categories:** Stores business categories.
- **sme_categories:** Stores many-to-many relation between SMEs and categories.
- **transactions:** Stores invoice, amount, payment method, and status.
- **payment_logs:** Stores payment gateway request and webhook history.

PostGIS should be enabled for location search:

```sql
CREATE EXTENSION IF NOT EXISTS postgis;
```

### Step 14: End-to-End MVP Flow

The MVP implementation should prove the complete business flow:

```text
SME registers account
   ↓
SME completes profile with category and lat/lng
   ↓
Corporation logs in
   ↓
Corporation searches SMEs by category and nearest location
   ↓
Corporation selects SME and creates transaction
   ↓
Payment virtual account is created
   ↓
Payment webhook marks transaction as paid
   ↓
RabbitMQ publishes payment.paid event
   ↓
WhatsApp notification is sent
   ↓
Report data is available
```

### Step 15: Suggested Development Order

Build the project in this order to reduce integration risk:

1. **Monorepo structure and shared packages**
2. **Database migrations and GORM models**
3. **Auth Service with JWT**
4. **API Gateway with JWT middleware**
5. **SME Service with category data**
6. **Nearby Service with PostGIS query**
7. **Transaction Service**
8. **Payment Gateway Service with Xendit webhook**
9. **RabbitMQ event publishing**
10. **Communication Service consumer**
11. **Report Service**
12. **OpenAPI/Swagger documentation**
