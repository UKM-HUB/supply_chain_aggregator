# ERD — Nearby Service

```mermaid
erDiagram
    SME_LOCATION {
        string  id          PK  "projected from smes.id"
        string  name            "projected from smes.name"
        string  address         "projected from smes.address"
        string  description     "projected from smes.description"
        array   category_ids    "projected from sme_categories"
        float   latitude
        float   longitude
        string  status          "active | inactive"
    }

    NEARBY_SME {
        string  id          PK  "embeds SMELocation"
        string  name
        string  address
        string  description
        array   category_ids
        float   latitude
        float   longitude
        string  status
        float   distance_km     "computed: ST_Distance result"
    }

    %% Each SMELocation is promoted to zero or one NearbySME per query (within radius → one result, outside → zero)
    %% Each NearbySME is derived from exactly one SMELocation
    SME_LOCATION |o--|| NEARBY_SME : "projected into when within radius"
```

## Cardinality rationale
| Relationship | Left | Right | Reason |
|---|---|---|---|
| SME_LOCATION → NEARBY_SME | zero or one | exactly one | An SME location either falls within the search radius (produces one result row) or does not (produces none). Each result always traces back to exactly one source location. |

## Notes
- This service owns **no persistent table**; it queries the `smes` table (sme-service) via PostGIS `ST_DWithin`.
- `SMELocation` and `NearbySME` are **read-only projections** (in-memory structs), not separate DB tables.
- `distance_km` is a computed field derived from `ST_Distance(location, user_point)`.
- Input: user latitude/longitude + optional `radius_km` and `category_id` query parameters.
- Output: ranked list of `NearbySME` sorted by `distance_km` ascending.
