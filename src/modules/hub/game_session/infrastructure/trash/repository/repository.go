package trash_repo

import (
	"context"
	"errors"

	"github.com/root9464/Go_GamlerDefi/src/packages/lib/logger"
	"gorm.io/gorm"
)

type TrashRepository struct {
	logger *logger.Logger
	db     *gorm.DB
}

type Ticket struct {
	GameID        string `json:"game_id"`
	GameSessionID uint   `json:"game_session_id"`
	PlayerID      string `json:"player_id"`
	Code          string `json:"code"`
}

func (Ticket) TableName() string {
	return "ticket"
}

type Session struct {
	HostID        string `json:"host_id"`
	GameSessionID uint   `json:"game_session_id"`
}

func (Session) TableName() string {
	return "game_session"
}

func NewTrashTicketRepository(logger *logger.Logger, db *gorm.DB) *TrashRepository {
	return &TrashRepository{
		logger: logger,
		db:     db,
	}
}

func (r *TrashRepository) IsUserHasAccess(ctx context.Context, sessionID uint, userID string) (bool, error) {
	r.logger.Infof("check access for session: %v and user: %s", sessionID, userID)

	var ticket Ticket
	err := r.db.WithContext(ctx).
		Where(&Ticket{
			GameSessionID: sessionID,
			PlayerID:      userID,
		}).
		First(&ticket).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.logger.Warnf("Access denied: No ticket found for user %s in session %v", userID, sessionID)
			return false, nil
		}
		r.logger.Errorf("Database error while checking access for user %s: %v", userID, err)
		return false, err
	}

	r.logger.Infof("Access granted: Ticket found for user %s in session %v", userID, sessionID)
	return true, nil
}

func (r *TrashRepository) GetSessionByID(ctx context.Context, sessionID uint) (*Session, error) {
	var session Session
	err := r.db.WithContext(ctx).
		Where(&Session{
			GameSessionID: sessionID,
		}).
		First(&session).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.logger.Warnf("Session not found for session %v", sessionID)
			return nil, nil
		}
		r.logger.Errorf("Database error while getting session for session %v: %v", sessionID, err)
		return nil, err
	}

	r.logger.Infof("Session found for session %v: %s", sessionID, session.HostID)
	return &session, nil
}

func (r *TrashRepository) GetSessionHost(ctx context.Context, sessionID uint) (string, error) {
	var session Session
	err := r.db.WithContext(ctx).
		Where(&Session{
			GameSessionID: sessionID,
		}).
		First(&session).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.logger.Warnf("Session host not found for session %v", sessionID)
			return "", nil
		}
		r.logger.Errorf("Database error while getting session host for session %v: %v", sessionID, err)
		return "", err
	}

	r.logger.Infof("Session host found for session %v: %s", sessionID, session.HostID)
	return session.HostID, nil
}
