# ERD — Report Service

```mermaid
erDiagram
    TRANSACTION_RECORD {
        string  id             PK  "projected from transactions.id"
        string  invoice_number     "projected from transactions.invoice_number"
        string  user_id            "projected from transactions.user_id"
        float   amount             "projected from transactions.amount"
        string  status             "pending | paid | failed | cancelled"
        time    created_at         "projected from transactions.created_at"
    }

    REPORT_SUMMARY {
        int   total_transaction "count of records in date range"
        float total_paid        "sum of amount WHERE status = paid"
        int   total_pending     "count WHERE status = pending"
    }

    CSV_EXPORT {
        string  id
        string  invoice_number
        string  user_id
        float   amount
        string  status
        time    created_at
    }

    %% Zero or many transaction records are aggregated into exactly one summary per report request
    TRANSACTION_RECORD }o--||  REPORT_SUMMARY : "aggregated into"

    %% Zero or many transaction records are serialized into exactly one CSV export per request
    TRANSACTION_RECORD }o--||  CSV_EXPORT     : "exported into"
```

## Cardinality rationale
| Relationship | Left | Right | Reason |
|---|---|---|---|
| TRANSACTION_RECORD → REPORT_SUMMARY | zero or many | exactly one | A date range may have no transactions (empty summary); all records in range collapse into one summary object |
| TRANSACTION_RECORD → CSV_EXPORT | zero or many | exactly one | A date range may export zero rows; all rows in range serialize into one CSV file |

## Notes
- This service owns **no persistent table**; it reads from the `transactions` table (transaction-service).
- `TransactionRecord`, `ReportSummary`, and `CSV_EXPORT` are **in-memory projections**, not DB tables.
- Three report endpoints:
  - `GET /api/v1/reports/daily?date=YYYY-MM-DD` — summary for a single day.
  - `GET /api/v1/reports/monthly?year=YYYY&month=MM` — summary for a full month.
  - `GET /api/v1/reports/export?from=YYYY-MM-DD&to=YYYY-MM-DD` — CSV download (`Content-Disposition: attachment`).
