# ERD — Payment Service

```mermaid
erDiagram
    TRANSACTIONS {
        VARCHAR(100) invoice_number PK "owned by transaction-service"
    }

    PAYMENT_LOGS {
        UUID          id               PK
        VARCHAR(100)  invoice_number   FK  "NOT NULL → TRANSACTIONS.invoice_number"
        NUMERIC(15_2) amount           "NOT NULL"
        VARCHAR(20)   user_phone
        TEXT          payment_url      "Xendit-hosted payment page URL"
        VARCHAR(255)  xendit_invoice_id "Xendit internal ID"
        VARCHAR(20)   status           "pending | paid | failed"
        TIMESTAMPTZ   created_at
        TIMESTAMPTZ   updated_at
    }

    XENDIT_INVOICE {
        string external_id   "= invoice_number (idempotency key)"
        float  amount
        string currency      "IDR"
        string description
        string invoice_url   "returned by Xendit API"
        string xendit_id     "Xendit internal ID"
        string status
    }

    RABBITMQ_EVENT {
        string invoice "invoice_number"
        float  amount
        string phone   "user_phone"
    }

    %% One transaction has zero or one payment log; each payment log covers exactly one transaction
    TRANSACTIONS    ||--o|  PAYMENT_LOGS    : "paid via (invoice_number)"

    %% Each payment log is created by calling exactly one Xendit invoice request
    PAYMENT_LOGS    ||--||  XENDIT_INVOICE  : "created via Xendit API"

    %% Each payment log that reaches 'paid' publishes exactly one event; each event traces to exactly one log
    PAYMENT_LOGS    ||--o|  RABBITMQ_EVENT  : "publishes on status=paid"
```

## Cardinality rationale
| Relationship | Left | Right | Reason |
|---|---|---|---|
| TRANSACTIONS → PAYMENT_LOGS | exactly one | zero or one | A transaction has no payment log until `create-va` is called; at most one log per transaction |
| PAYMENT_LOGS → XENDIT_INVOICE | exactly one | exactly one | Creating a payment log always involves one Xendit API call (real or mock) |
| PAYMENT_LOGS → RABBITMQ_EVENT | exactly one | zero or one | Only `paid` webhooks publish an event; `pending`/`failed` logs never emit one |

## Notes
- `XENDIT_INVOICE` and `RABBITMQ_EVENT` are **not DB tables** — they represent the Xendit API payload and the RabbitMQ message body.
- `invoice_number` is a cross-service FK to `transactions.invoice_number` (transaction-service), enforced in the DB migration.
- When `XENDIT_SECRET_KEY` env var is empty the client runs in **mock mode** (returns a fake invoice URL).
