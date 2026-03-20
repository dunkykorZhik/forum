package repository

import (
	"database/sql"
	"fmt"
	"forum/internal/model"
	"strings"
	"time"
)

type SessionRepo struct {
	db *sql.DB
}

func NewSessionRepo(db *sql.DB) *SessionRepo {
	return &SessionRepo{db}
}

func (s *SessionRepo) GetByUuid(uuid string) (*model.Session, error) {
	row := s.db.QueryRow(`
SELECT id, uuid, expired_at, user_id FROM sessions
WHERE uuid = ?`, uuid)
	session := &model.Session{}
	strExpiredAt := ""

	err := row.Scan(&session.Id, &session.Uuid, &strExpiredAt, &session.UserId)

	switch {
	case err == nil:
		timeExpiredAt, err := time.ParseInLocation(time.RFC3339, strExpiredAt, time.Local)
		if err != nil {
			return nil, fmt.Errorf("time.Parse: %w", err)
		}
		session.ExpiredAt = timeExpiredAt
		return session, nil
	case strings.HasPrefix(err.Error(), "sql: no rows in result set"):
		return nil, model.ErrSessNotFound
	default:
		return nil, fmt.Errorf("row.Scan: %w", err)
	}
}
func (s *SessionRepo) UpdateByUserId(userId int64, session *model.Session) error {
	strExpiredAt := session.ExpiredAt.Format(time.RFC3339)
	row := s.db.QueryRow(`
UPDATE sessions 
SET uuid = ?, expired_at = ?
WHERE user_id = ?
RETURNING id`, session.Uuid, strExpiredAt, session.UserId)

	err := row.Scan(&session.Id)
	switch {
	case err == nil:
		return nil
	}
	return fmt.Errorf("row.Scan: %w", err)
}
func (s *SessionRepo) Create(session *model.Session) (int64, error) {
	strExpiredAt := session.ExpiredAt.Format(time.RFC3339)
	row := s.db.QueryRow(`
INSERT INTO sessions (uuid, expired_at, user_id) VALUES
(?, ?, ?) RETURNING id`, session.Uuid, strExpiredAt, session.UserId)

	err := row.Scan(&session.Id)
	switch {
	case err == nil:
		return session.Id, nil
	case strings.HasPrefix(err.Error(), "constraint failed: UNIQUE constraint failed"):
		return -1, model.ErrSessionExists
	}
	return -1, fmt.Errorf("row.Scan: %w", err)
}
