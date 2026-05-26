# ERD — SME Service

```mermaid
erDiagram
    USERS {
        UUID id PK "owned by auth-service"
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
        TEXT[]       products    "array of product names"
        VARCHAR(255) capacity
        DOUBLE       latitude    "NOT NULL"
        DOUBLE       longitude   "NOT NULL"
        GEOGRAPHY    location    "POINT 4326, auto-synced via trigger"
        VARCHAR(20)  status      "active | inactive"
        TIMESTAMPTZ  created_at
        TIMESTAMPTZ  updated_at
    }

    SME_CATEGORIES {
        UUID         sme_id      PK,FK "NOT NULL → SMES.id"
        VARCHAR(100) category_id PK,FK "NOT NULL → CATEGORIES.id"
    }

    %% One user owns zero or many SMEs; each SME belongs to exactly one user
    USERS           ||--o{  SMES           : "owns (owner_id)"

    %% One SME has zero or many category assignments; each assignment belongs to exactly one SME
    SMES            ||--o{  SME_CATEGORIES : "assigned to"

    %% One category appears in zero or many assignments; each assignment points to exactly one category
    CATEGORIES      ||--o{  SME_CATEGORIES : "tagged on"
```

## Cardinality rationale
| Relationship | Left | Right | Reason |
|---|---|---|---|
| USERS → SMES | exactly one | zero or many | A user may own no SME yet; they can register multiple |
| SMES → SME_CATEGORIES | exactly one | zero or many | Category tags are optional at creation time |
| CATEGORIES → SME_CATEGORIES | exactly one | zero or many | A category may exist before any SME is tagged with it |

## Notes
- `owner_id` is a cross-service FK to `users.id` (auth-service), enforced in the DB migration.
- `location` (PostGIS GEOGRAPHY) is auto-populated from `latitude`/`longitude` via the `trg_sme_location` trigger on INSERT/UPDATE.
- `products` is a PostgreSQL text array; each element is a product/service name offered by the SME.
- GIST index on `location` enables fast `ST_DWithin` queries used by nearby-service.
