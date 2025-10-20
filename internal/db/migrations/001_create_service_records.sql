CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS service_scans (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    ip VARCHAR(45) NOT NULL,
    port INTEGER NOT NULL,
    service VARCHAR(50) NOT NULL,
    response TEXT NOT NULL,
    last_scanned TIMESTAMP NOT NULL,
    UNIQUE(ip, port, service)
);

CREATE INDEX IF NOT EXISTS idx_service_scans_lookup ON service_scans(ip, port, service);
CREATE INDEX IF NOT EXISTS idx_service_scans_last_scanned ON service_scans(last_scanned);
