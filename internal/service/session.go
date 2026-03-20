package service

import (
	"errors"
	"fmt"
	"forum/internal/model"
	"forum/internal/repository"
	"time"

	"github.com/google/uuid"
)

type SessionService struct {
	repo repository.Session
}

func NewSessionService(repo repository.Session) *SessionService {
	return &SessionService{repo}
}

func (s *SessionService) GetByUuid(uuid string) (*model.Session, error) {
	session, err := s.repo.GetByUuid(uuid)
	switch {
	case err == nil:
		expiredInSec := time.Until(session.ExpiredAt).Seconds()
		if expiredInSec <= 0 {
			return nil, model.ErrExpired
		}
		return session, nil
	case errors.Is(err, model.ErrNotFound):
		return nil, model.ErrSessNotFound
	default:
		return nil, fmt.Errorf("s.repo.GetByUuid: %w", err)
	}
}

func (s *SessionService) Record(userId int64) (*model.Session, error) {
	uid := uuid.New()
	session := &model.Session{
		Uuid:      uid.String(),
		UserId:    userId,
		ExpiredAt: time.Now().Add(time.Minute * 15),
	}

	_, err := s.repo.Create(session)
	switch {
	case err == nil:
		return session, nil
	case errors.Is(err, model.ErrSessionExists):
		err := s.repo.UpdateByUserId(session.UserId, session)
		if err != nil {
			return nil, fmt.Errorf("s.repo.UpdateByUserId: %w", err)
		}
		return session, nil
	default:
		return nil, fmt.Errorf("s.repo.Create: %w", err)
	}
}
