package service

import (
	"errors"
	"fmt"
	"forum/internal/model"
	"forum/internal/repository"
	"strings"
	"time"
)

type UserService struct {
	repo repository.User
}

func NewUserService(repo repository.User) *UserService {
	return &UserService{repo}
}

func (u *UserService) Create(user *model.User) (int64, error) {
	if err := user.ValidateNickname(); err != nil {
		return -1, model.ErrInvalidUsername
	} else if err := user.ValidateEmail(); err != nil {
		return -1, model.ErrInvalidEmail
	}

	err := user.HashPassword()
	if err != nil {
		return -1, fmt.Errorf("user.HashPassword: %w", err)
	}

	user.CreatedAt = time.Now()
	userId, err := u.repo.Create(user)

	switch {
	case err == nil:
		return userId, nil
	case errors.Is(err, model.ErrExistEmail),
		errors.Is(err, model.ErrExistUsername),
		errors.Is(err, model.ErrWrongLengthEmail),
		errors.Is(err, model.ErrWrongLengthUsername):
		return -1, err
	}
	return -1, fmt.Errorf("u.repo.Create: %w", err)
}
func (u *UserService) GetByID(id int64) (*model.User, error) {
	usr, err := u.repo.GetByID(id)
	switch {
	case err == nil:
		return usr, nil
	case errors.Is(err, model.ErrNotFound):
		return nil, model.ErrNotFound
	}
	return nil, fmt.Errorf("u.repo.GetByID: %w", err)
}

func (u *UserService) GetByUsernameOrEmail(field string) (*model.User, error) {
	switch {
	case strings.Contains(field, "@"):
		if err := (&model.User{Email: field}).ValidateEmail(); err != nil {
			return nil, model.ErrInvalidEmail
		}
		usr, err := u.repo.GetByEmail(field)
		switch {
		case err == nil:
			return usr, err
		case errors.Is(err, model.ErrNotFound):
			return nil, model.ErrNotFound
		}
		return nil, fmt.Errorf("u.repo.GetByEmail: %w", err)
	default:
		if err := (&model.User{Username: field}).ValidateNickname(); err != nil {
			return nil, model.ErrInvalidUsername
		}
		usr, err := u.repo.GetByUsername(field)
		switch {
		case err == nil:
			return usr, err
		case errors.Is(err, model.ErrNotFound):
			return nil, err
		}
		return nil, fmt.Errorf("u.repo.GetByNickname: %w", err)
	}
}
