package repository

import (
	"database/sql"
	"fmt"
	"forum/internal/model"
	"strings"
	"time"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db}
}
func (u *UserRepo) GetByUsername(nickname string) (*model.User, error) {
	row := u.db.QueryRow(`
SELECT id, nickname, email, password, created_at FROM users
WHERE nickname = ?`, nickname)
	user := &model.User{}
	strCreatedAt := ""

	err := row.Scan(&user.Id, &user.Username, &user.Email, &user.Password, &strCreatedAt)

	switch {
	case err == nil:
		timeCreatedAt, err := time.ParseInLocation(time.RFC3339, strCreatedAt, time.Local)
		if err != nil {
			return nil, fmt.Errorf("time.Parse: %w", err)
		}
		user.CreatedAt = timeCreatedAt
		return user, nil
	case strings.HasPrefix(err.Error(), "sql: no rows in result set"):
		return nil, model.ErrNotFound
	default:
		return nil, fmt.Errorf("row.Scan: %w", err)
	}
}

func (u *UserRepo) GetByID(id int64) (*model.User, error) {
	row := u.db.QueryRow(`
SELECT id, nickname, email, created_at FROM users
WHERE id = ?`, id)
	user := &model.User{}
	strCreatedAt := ""
	err := row.Scan(&user.Id, &user.Username, &user.Email, &strCreatedAt)
	switch {
	case err == nil:
		timeCreatedAt, err := time.ParseInLocation(time.RFC3339, strCreatedAt, time.Local)
		if err != nil {
			return nil, fmt.Errorf("time.Parse: %w", err)
		}
		user.CreatedAt = timeCreatedAt
		return user, nil
	case strings.HasPrefix(err.Error(), "sql: no rows in result set"):
		return nil, model.ErrNotFound
	default:
		return nil, fmt.Errorf("row.Scan: %w", err)
	}
}
func (u *UserRepo) GetByEmail(email string) (*model.User, error) {
	row := u.db.QueryRow(`
SELECT id, nickname, email, password, created_at FROM users
WHERE email = ?`, email)
	user := &model.User{}
	strCreatedAt := ""

	err := row.Scan(&user.Id, &user.Username, &user.Email, &user.Password, &strCreatedAt)

	switch {
	case err == nil:
		timeCreatedAt, err := time.ParseInLocation(time.RFC3339, strCreatedAt, time.Local)
		if err != nil {
			return nil, fmt.Errorf("time.Parse: %w", err)
		}
		user.CreatedAt = timeCreatedAt
		return user, nil
	case strings.HasPrefix(err.Error(), "sql: no rows in result set"):
		return nil, model.ErrNotFound
	default:
		return nil, fmt.Errorf("row.Scan: %w", err)
	}
}

func (u *UserRepo) Create(user *model.User) (int64, error) {
	strCreatedAt := user.CreatedAt.Format(time.RFC3339)
	row := u.db.QueryRow(`
INSERT INTO users (nickname, email, password, created_at) VALUES
(?, ?, ?, ?) RETURNING id`, user.Username, user.Email, user.Password, strCreatedAt)

	err := row.Scan(&user.Id)
	switch {
	case err == nil:
		return user.Id, nil
	case strings.HasPrefix(err.Error(), "constraint failed: UNIQUE constraint failed"):
		switch {
		case strings.Contains(err.Error(), "nickname"):
			return -1, model.ErrExistUsername
		case strings.Contains(err.Error(), "email"):
			return -1, model.ErrExistEmail
		}
	case strings.HasPrefix(err.Error(), "constraint failed: CHECK constraint failed"):
		switch {
		case strings.Contains(err.Error(), "LENGTH(nickname)"):
			return -1, model.ErrWrongLengthUsername
		case strings.Contains(err.Error(), "LENGTH(email)"):
			return -1, model.ErrWrongLengthEmail
		}
	}
	return -1, fmt.Errorf("row.Scan: %w", err)
}
