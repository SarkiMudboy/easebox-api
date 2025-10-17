-- Drop views
-- DROP VIEW IF EXISTS session_statistics;
-- DROP VIEW IF EXISTS active_session_latest_locations;

-- Drop triggers
DROP TRIGGER IF EXISTS trigger_update_tracking_sessions_updated_at ON tracking_sessions;

-- Drop functions
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop tables (cascade will remove foreign key constraints)
DROP TABLE IF EXISTS location_updates CASCADE;
DROP TABLE IF EXISTS tracking_sessions CASCADE;
