# ERD — User Service

```mermaid
erDiagram
    USERS {
        UUID         id         PK  "projected from auth-service users table"
        VARCHAR(255) name
        VARCHAR(255) email
        VARCHAR(20)  phone
        VARCHAR(20)  role       "ADMIN | CORPORATION | SME"
        VARCHAR(20)  status     "active | inactive"
        TIMESTAMPTZ  created_at
        TIMESTAMPTZ  updated_at
    }

    SMES {
        UUID owner_id FK "references USERS.id"
    }

    TRANSACTIONS {
        UUID user_id FK "references USERS.id"
    }

    %% One user owns zero or many SMEs; each SME belongs to exactly one user
    USERS ||--o{ SMES         : "owns (owner_id)"

    %% One user places zero or many transactions; each transaction belongs to exactly one user
    USERS ||--o{ TRANSACTIONS : "places (user_id)"
```

## Cardinality rationale
| Relationship | Left | Right | Reason |
|---|---|---|---|
| USERS → SMES | exactly one | zero or many | A user may own no SME yet; a user with role SME can register multiple business profiles |
| USERS → TRANSACTIONS | exactly one | zero or many | A user may have no transactions yet; they can accumulate many over time |

## Notes
- This service is a **read-only view** over the `users` table owned by auth-service.
- It does not write to any table; mutations (register, login) go through auth-service.
- `password_hash` is **never returned** by this service.
- gRPC interface defined in `proto/user/user.proto` (`GetUser`, `ListUsers`).
