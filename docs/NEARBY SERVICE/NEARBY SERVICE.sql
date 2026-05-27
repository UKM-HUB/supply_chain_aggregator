CREATE TABLE umkms (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    address TEXT,
    location GEOGRAPHY(POINT, 4326)
);

//Insert sample data

INSERT INTO umkms (name, address, location)
VALUES
(
    'Warung Makan Berkah',
    'Jakarta Pusat',
    ST_SetSRID(ST_MakePoint(106.816666, -6.200000), 4326)
),
(
    'Kopi Nusantara',
    'Jakarta Selatan',
    ST_SetSRID(ST_MakePoint(106.827153, -6.175110), 4326)
),
(
    'Bakso Enak',
    'Jakarta Barat',
    ST_SetSRID(ST_MakePoint(106.790000, -6.210000), 4326)
);