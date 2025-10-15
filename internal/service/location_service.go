package service

import (
	"context"
	"fmt"
	"time"

	"github.com/SarkiMudboy/easebox-api/internal/domain"
	"github.com/SarkiMudboy/easebox-api/internal/repository"
)

type LocationService struct {
	locationRepo repository.LocationRepository
	sessionRepo repository.SessionRepository
}

func NewLocationService(locationRepo repository.LocationRepository, sessionRepo repository.SessionRepository) *LocationService {
	return &LocationService{
		locationRepo: locationRepo,
		sessionRepo: sessionRepo,
	}
}

func (s *LocationService) RecordLocation(ctx context.Context, location *domain.LocationUpdate) error {
	if err := s.validateLocation(location); err != nil {
		return err
	}

	session, err := s.sessionRepo.GetByID(ctx, location.SessionID)
	if err != nil {
		return fmt.Errorf("session not found: %w", err)
	}

	if !session.IsActive {
		return &domain.DomainError{
			Code: "SESSION_INACTIVE",
			Message: "cannot record location for inactive session",
		}
	}

	if err := s.locationRepo.Create(ctx, location); err != nil {
		return fmt.Errorf("failed to record location: %w", err)
	}

	return nil
}

func (s *LocationService) StartTracking(ctx context.Context, sessionID, deliveryID string) error {
	session := &domain.TrackingSession{
		SessionID: sessionID,
		DeliveryID: deliveryID,
		StartTime: time.Now(),
		IsActive: true,
	}

	return s.sessionRepo.Create(ctx, session)
}

func (s *LocationService) StopTracking(ctx context.Context, sessionID string) error {
	session, err := s.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return err
	}

	now := time.Now()
	session.EndTime = &now
	session.IsActive = false

	return s.sessionRepo.Update(ctx, session)
}

func (s *LocationService) GetSessionRoute(ctx context.Context, sessionID string) ([]*domain.LocationUpdate, error) {
	return s.locationRepo.GetBySessionID(ctx, sessionID)
}

func (s *LocationService) validateLocation (location *domain.LocationUpdate) error {
	if location.Latitude < -90 || location.Latitude > 90 {
		return &domain.DomainError{Code: "INVALID_LATITUDE", Message: "Latitude must be between -90 and 90"}
	}
	if location.Longitude < -180 || location.Longitude > 180 {
		return &domain.DomainError{Code: "INVALID_LONGITUDE", Message: "Longitude must be between -180 and 180"}
	}

	if location.Accuracy < 0 {
		return &domain.DomainError{Code: "INVALID_ACCURACY", Message: "accuracy cannot be negative"}
	}

	return nil
}