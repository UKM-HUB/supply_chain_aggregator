# ERD — Auth Service

```mermaid
erDiagram
    USERS {
        UUID         id            PK
        VARCHAR(255) name          "NOT NULL"
        VARCHAR(255) email         "NOT NULL, UNIQUE"
        VARCHAR(20)  phone
        VARCHAR(255) password_hash "NOT NULL"
        VARCHAR(20)  role          "ADMIN | CORPORATION | SME"
        VARCHAR(20)  status        "active | inactive"
        DOUBLE       latitude
        DOUBLE       longitude
        TIMESTAMPTZ  created_at
        TIMESTAMPTZ  updated_at
    }

    SMES {
        UUID owner_id FK "references USERS.id"
    }

    TRANSACTIONS {
        UUID user_id FK "references USERS.id"
    }

    USERS ||--o{ SMES         : "one user owns zero or many SMEs"
    USERS ||--o{ TRANSACTIONS : "one user places zero or many transactions"
```

## Cardinality legend
| Notation | Meaning |
|----------|---------|
| `\|\|`  | exactly one |
| `o\|`   | zero or one |
| `\|{`   | one or more |
| `o{`    | zero or more |

## Notes
- Single source of truth for user identity across all services.
- Issues JWT tokens; downstream services verify the token independently.
- `role` drives authorization: `CORPORATION` creates transactions, `SME` registers business profiles, `ADMIN` has full access.
- `latitude` / `longitude` store the user's location, used by nearby-service to find SMEs near the requester.
- Every `smes.owner_id` and `transactions.user_id` is a cross-service FK back to `USERS.id`.
