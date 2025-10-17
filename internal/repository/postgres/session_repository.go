package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/SarkiMudboy/easebox-api/internal/domain"
	"github.com/SarkiMudboy/easebox-api/internal/repository"
)

type SessionRepository struct {
	db *sql.DB
}

func NewSessionRepository(db *sql.DB) repository.SessionRepository {
	return &SessionRepository{db: db}
}

func (r *SessionRepository) Create(ctx context.Context, session *domain.TrackingSession) error {

	query := `
		INSERT INTO tracking_sessions
			(session_id, delivery_id, start_time, is_active)
		VALUES ($1, $2, $3, $4);
	`

	_, err := r.db.ExecContext(ctx, query, session.SessionID, session.DeliveryID, session.StartTime, session.IsActive)

	if err != nil {
		return fmt.Errorf("Failed to create session for %v: %v", session.SessionID, err)
	}

	return nil
}

func (r *SessionRepository) GetByID(ctx context.Context, sessionID string) (session *domain.TrackingSession, err error) {

	session = &domain.TrackingSession{
		SessionID: sessionID,
	}

	query := `
		SELECT 
			delivery_id, start_time, end_time, is_active
		FROM tracking_sessions 
		WHERE session_id = $1;
	`

	err = r.db.QueryRowContext(ctx, query, sessionID).Scan(&session.DeliveryID, &session.StartTime, &session.EndTime, &session.IsActive)

	return
}

func (r *SessionRepository) Update(ctx context.Context, session *domain.TrackingSession) error {
	query := `
		UPDATE tracking_sessions 
			SET delivery_id = $2, start_time = $3, end_time = $4, is_active = $5
		WHERE session_id = $1;
	`

	_, err := r.db.ExecContext(ctx, query, session.SessionID, session.DeliveryID, session.StartTime, session.EndTime, session.IsActive)
	if err != nil {
		return fmt.Errorf("Failed to update tracking session %v: %v", session.SessionID, err)
	}

	return nil
}