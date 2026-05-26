CREATE TABLE smes (
    id          UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_id    UUID         NOT NULL REFERENCES users (id) ON DELETE RESTRICT,
    name        VARCHAR(255) NOT NULL,
    phone       VARCHAR(20),
    address     TEXT,
    description TEXT,
    products    TEXT[]       NOT NULL DEFAULT '{}',
    capacity    VARCHAR(255),
    latitude    DOUBLE PRECISION NOT NULL,
    longitude   DOUBLE PRECISION NOT NULL,
    location    GEOGRAPHY(POINT, 4326),
    status      VARCHAR(20)  NOT NULL DEFAULT 'active',
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    CONSTRAINT smes_status_check CHECK (status IN ('active', 'inactive'))
);

-- Auto-populate PostGIS location column from latitude/longitude
CREATE OR REPLACE FUNCTION sync_sme_location()
RETURNS TRIGGER AS $$
BEGIN
    NEW.location = ST_SetSRID(ST_MakePoint(NEW.longitude, NEW.latitude), 4326)::geography;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_sme_location
    BEFORE INSERT OR UPDATE OF latitude, longitude ON smes
    FOR EACH ROW EXECUTE FUNCTION sync_sme_location();

CREATE INDEX idx_smes_owner_id ON smes (owner_id);
CREATE INDEX idx_smes_status   ON smes (status);
CREATE INDEX idx_smes_location ON smes USING GIST (location);
