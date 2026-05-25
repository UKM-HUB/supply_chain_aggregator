PostgreSQL PostGIS
PostGIS Extension

Endpoint
GET /api/v1/nearby/umkm?lat=-6.2&lng=106.8

Query 
SELECT *,
ST_Distance(
location,
ST_MakePoint(106.8, -6.2)
)
FROM umkms
ORDER BY location <-> ST_MakePoint(106.8, -6.2)
LIMIT 10;