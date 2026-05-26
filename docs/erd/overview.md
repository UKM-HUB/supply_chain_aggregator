# ERD — Cross-Service Overview

All persistent tables and their relationships across service boundaries.

```mermaid
erDiagram
    USERS {
        UUID         id            PK
        VARCHAR(255) name
        VARCHAR(255) email         "UNIQUE"
        VARCHAR(20)  phone
        VARCHAR(255) password_hash
        VARCHAR(20)  role          "ADMIN | CORPORATION | SME"
        VARCHAR(20)  status        "active | inactive"
        DOUBLE       latitude
        DOUBLE       longitude
        TIMESTAMPTZ  created_at
        TIMESTAMPTZ  updated_at
    }

    CATEGORIES {
        VARCHAR(100) id          PK
        VARCHAR(255) name        "NOT NULL"
        TEXT         description
        TIMESTAMPTZ  created_at
    }

    SMES {
        UUID         id          PK
        UUID         owner_id    FK  "NOT NULL → USERS.id"
        VARCHAR(255) name        "NOT NULL"
        VARCHAR(20)  phone
        TEXT         address
        TEXT         description
        TEXT[]       products
        VARCHAR(255) capacity
        DOUBLE       latitude    "NOT NULL"
        DOUBLE       longitude   "NOT NULL"
        GEOGRAPHY    location    "POINT 4326, auto via trigger"
        VARCHAR(20)  status      "active | inactive"
        TIMESTAMPTZ  created_at
        TIMESTAMPTZ  updated_at
    }

    SME_CATEGORIES {
        UUID         sme_id      PK,FK "NOT NULL → SMES.id"
        VARCHAR(100) category_id PK,FK "NOT NULL → CATEGORIES.id"
    }

    TRANSACTIONS {
        UUID          id             PK
        VARCHAR(100)  invoice_number "NOT NULL, UNIQUE"
        UUID          user_id        FK  "NOT NULL → USERS.id"
        NUMERIC(15_2) amount         "NOT NULL, > 0"
        VARCHAR(20)   status         "pending | paid | failed | cancelled"
        VARCHAR(50)   payment_method
        TIMESTAMPTZ   created_at
        TIMESTAMPTZ   updated_at
    }

    PAYMENT_LOGS {
        UUID          id               PK
        VARCHAR(100)  invoice_number   FK  "NOT NULL → TRANSACTIONS.invoice_number"
        NUMERIC(15_2) amount           "NOT NULL"
        VARCHAR(20)   user_phone
        TEXT          payment_url
        VARCHAR(255)  xendit_invoice_id
        VARCHAR(20)   status           "pending | paid | failed"
        TIMESTAMPTZ   created_at
        TIMESTAMPTZ   updated_at
    }

    %% One user owns zero or many SMEs; each SME belongs to exactly one user
    USERS        ||--o{  SMES            : "owns (owner_id)"

    %% One user places zero or many transactions; each transaction belongs to exactly one user
    USERS        ||--o{  TRANSACTIONS    : "places (user_id)"

    %% One SME has zero or many category assignments; each assignment belongs to exactly one SME
    SMES         ||--o{  SME_CATEGORIES  : "assigned to"

    %% One category appears in zero or many assignments; each assignment points to exactly one category
    CATEGORIES   ||--o{  SME_CATEGORIES  : "tagged on"

    %% One transaction has zero or one payment log; each payment log covers exactly one transaction
    TRANSACTIONS ||--o|  PAYMENT_LOGS    : "paid via (invoice_number)"
```

## Cardinality reference

| Relationship | Notation | Left side | Right side |
|---|---|---|---|
| USERS → SMES | `\|\|--o{` | exactly one user | zero or many SMEs |
| USERS → TRANSACTIONS | `\|\|--o{` | exactly one user | zero or many transactions |
| SMES → SME_CATEGORIES | `\|\|--o{` | exactly one SME | zero or many category rows |
| CATEGORIES → SME_CATEGORIES | `\|\|--o{` | exactly one category | zero or many SME rows |
| TRANSACTIONS → PAYMENT_LOGS | `\|\|--o\|` | exactly one transaction | zero or one payment log |

## Service ownership

| Table | Owner Service | Consumers |
|---|---|---|
| `users` | auth-service | user-service (read), nearby-service |
| `categories` | sme-service | nearby-service (read) |
| `smes` | sme-service | nearby-service (PostGIS read) |
| `sme_categories` | sme-service | nearby-service (read) |
| `transactions` | transaction-service | report-service (read), payment-service |
| `payment_logs` | payment-service | — |

## Services with no persistent table

| Service | Role |
|---|---|
| api-gateway | HTTP reverse proxy + JWT validation |
| nearby-service | PostGIS read queries on `smes`, returns `NearbySME` projections |
| user-service | Read-only projection of `users` |
| report-service | Read-only projection of `transactions` |
| communication-service | RabbitMQ consumer, calls WhatsApp API |
