package domain

import "time"

type LocationUpdate struct {
	ID         int64
	SessionID  string
	DeliveryID string
	Latitude   float64
	Longitude  float64
	Accuracy   float64
	Speed      *int64
	Heading    *int64
	RecordedAt time.Time
	CreatedAt time.Time 
}


type TrackingSession struct {
	SessionID string
	DeliveryID string
	StartTime time.Time
	EndTime *time.Time
	IsActive bool
}

type DomainError struct {
	Code string
	Message string 
}

func (e *DomainError) Error () string {
	return e.Message
}