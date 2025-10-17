-- Enable POSTGIS extension
CREATE EXTENSION IF NOT EXISTS postgis;

-- Enable UUID extension for generating UUIDs
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ============================================
-- Tracking Sessions Table
-- ============================================
-- Stores metadata about each tracking session
CREATE TABLE tracking_sessions (
    id BIGSERIAL PRIMARY KEY,
    session_id VARCHAR(255) NOT NULL UNIQUE,
    delivery_id UUID,
    is_active BOOLEAN NOT NULL DEFAULT true,
    start_time TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    end_time TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Constraints
    CONSTRAINT check_end_time_after_start CHECK (end_time IS null OR end_time >= start_time)  
);


-- Indices for tracking_sessions table
CREATE INDEX idx_tracking_sessions_session_id 
    ON tracking_sessions(session_id);
CREATE INDEX idx_tracking_sessions_delivery_id 
    ON tracking_sessions(delivery_id);
CREATE INDEX idx_tracking_sessions_is_active 
    ON tracking_sessions(is_active) WHERE is_active = true;
CREATE INDEX idx_tracking_sessions_start_time 
    ON tracking_sessions(start_time DESC);

-- ============================================
-- Location Updates Table
-- ============================================
-- Stores individual location points

CREATE TABLE location_updates (
    
    id BIGSERIAL PRIMARY KEY,
    session_id VARCHAR(255) NOT NULL,
    delivery_id UUID,

     -- SRID 4326 = WGS84 coordinate system (standard GPS coordinates)
    location GEOGRAPHY(POINT, 4326) NOT NULL,

    -- location metadata
    accuracy DOUBLE PRECISION NOT NULL CHECK (accuracy >= 0),
    speed DOUBLE PRECISION CHECK (speed IS NULL OR speed >= 0),
    heading DOUBLE PRECISION CHECK (heading IS NULL OR heading >= 0 AND heading <= 360),

    recorded_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_location_tracking_session FOREIGN KEY (session_id) 
        REFERENCES tracking_sessions(session_id)
        ON DELETE CASCADE

);

-- Indices for location_updates table 
CREATE INDEX idx_location_updates_session_time 
    ON location_updates(session_id, recorded_at DESC);

CREATE INDEX idx_location_updates_delivery_time
    ON location_updates(delivery_id, recorded_at DESC);

-- Spatial index for geographic queries (e.g. "find locations within X meters")
CREATE INDEX idx_location_updates_geography
    ON location_updates USING GIST(location);

-- Composite index for active tracking queries
CREATE INDEX idx_location_session_created
    ON location_updates(session_id, created_at DESC);

-- Index for time-based queries
CREATE INDEX idx_location_updates_recorded_at
    ON location_updates(recorded_at DESC);


-- helper functions (utils)

-- -- Function to automatically update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;

$$ LANGUAGE plpgsql;

-- Trigger to auto-update updated_at on tracking_sessions
CREATE TRIGGER trigger_update_tracking_sessions_updated_at
    BEFORE UPDATE ON tracking_sessions
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();