package repository

import (
	"context"

	"github.com/SarkiMudboy/easebox-api/internal/domain"
)



type LocationRepository interface {
	Create(ctx context.Context, location *domain.LocationUpdate) error
	GetBySessionID(ctx context.Context, sessionID string) ([]*domain.LocationUpdate, error)
	GetByDeliveryID(ctx context.Context, deliveryID string) ([]*domain.LocationUpdate, error)
	GetLatestBySessionID(ctx context.Context, sessionID string) (*domain.LocationUpdate, error)
	GetWithinRadius(ctx context.Context, lat, long, radiusMeters float64) ([]*domain.LocationUpdate, error)
}


type SessionRepository interface {
	Create(ctx context.Context, session *domain.TrackingSession) error
	GetByID(ctx context.Context, sessionID string) (*domain.TrackingSession, error)
	Update(ctx context.Context, session *domain.TrackingSession) error
}