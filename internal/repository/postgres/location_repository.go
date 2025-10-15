package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/SarkiMudboy/easebox-api/internal/domain"
	"github.com/SarkiMudboy/easebox-api/internal/repository"
)

type locationRepository struct {
	db *sql.DB
}

func NewLocationRepository(db *sql.DB) repository.LocationRepository {
	return &locationRepository{db: db}
}


func (r *locationRepository) Create(ctx context.Context, location *domain.LocationUpdate) error {
	query := `
		INSERT INTO location_updates 
			(session_id, delivery_id, location, accuracy, speed, heading, recorded_at) 
		VALUES 
			($1, $2, $3, ST_SetSRID(ST_MakePoint($3, $4), 4326), $5, $6, $7, $8) 
		RETURNING id, created_at
	`

	err := r.db.QueryRowContext(
		ctx, 
		query,
		location.SessionID, 
		location.DeliveryID, 
		location.Longitude,
		location.Latitude, 
		location.Accuracy, 
		location.Speed, 
		location.Heading, 
		location.RecordedAt,
		).Scan(
			&location.ID, &location.CreatedAt,
		)
	
	if err != nil {
		return fmt.Errorf("Failed to create location update entry: %w", err)
	}

	return nil
}

func (r *locationRepository) GetBySessionID(ctx context.Context, sessionID string) (locations []*domain.LocationUpdate, err error) {
	query := `
		SELECT 
			(id, delivery_id, ST_Y(location::geometry) as latitude, ST_X(location::geometry) as longitude, accuracy, speed, heading, recorded_at, created_at) 
		FROM location_updates 
			WHERE session_id = $1
		ORDER BY recorded_at ASC
	`
	rows, err := r.db.QueryContext(ctx, query, sessionID)
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve location update: %w", err)
	}

	for rows.Next() {
		location := &domain.LocationUpdate{}

		err = rows.Scan(
			&location.ID, location.DeliveryID, &location.Latitude, &location.Accuracy, &location.Speed, &location.Heading, &location.RecordedAt, &location.CreatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("Failed to Scan location: %w", err)
		}
		locations = append(locations, location)
	}

	defer rows.Close()

	return
}

func (r *locationRepository) GetByDeliveryID(ctx context.Context, deliveryID string) (locations []*domain.LocationUpdate, err error) {
	query := `
		SELECT 
			(id, session_id, ST_Y(location::geometry) as latitude, ST_X(location::geometry) as longitude, accuracy, speed, heading, recorded_at, created_at) 
		FROM location_updates 
			WHERE delivery_id = $1
		ORDER BY recorded_at ASC
	`
	rows, err := r.db.QueryContext(ctx, query, deliveryID)
	
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve location update: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		location := &domain.LocationUpdate{}

		err = rows.Scan(
			&location.ID, location.SessionID, &location.Latitude, &location.Accuracy, &location.Speed, &location.Heading, &location.RecordedAt, &location.CreatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("Failed to Scan location: %w", err)
		}
		locations = append(locations, location)
	}

	return
}

func (r *locationRepository) GetLatestBySessionID(ctx context.Context, sessionID string) (loc *domain.LocationUpdate, err error) {

	query := `
		SELECT DISTINCT ON (session_id)
			id,
			session_id,
			ST_Y(location::geometry) AS latitude,
			ST_X(location::geometry) AS longitude,
			accuracy,
			speed,
			heading,
			recorded_at,
			created_at
		FROM location_updates
		WHERE session_id = $1
		ORDER BY session_id, recorded_at DESC
		LIMIT 1;
	`

	err = r.db.QueryRowContext(ctx, query, sessionID).Scan(&loc)

	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve location: %v", err)
	}

	return
}

// func GetWithinRadius(ctx context.Context, lat, long, radiusMeters float64) ([]*domain.LocationUpdate, error) {}