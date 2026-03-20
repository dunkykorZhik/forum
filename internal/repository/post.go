package repository

import (
	"database/sql"
	"fmt"
	"forum/internal/model"
	"strings"
	"time"
)

type PostRepo struct {
	db *sql.DB
}

func NewPostRepo(db *sql.DB) *PostRepo {
	return &PostRepo{db}
}

func (p *PostRepo) GetByUserID(userId, offset, limit int64) ([]*model.Post, error) {
	if limit == 0 {
		limit = -1
	}

	rows, err := p.db.Query(`
SELECT id, title, content, user_id, created_at, updated_at FROM posts
WHERE user_id = ?
LIMIT ? OFFSET ?
	`, userId, limit, offset)

	if err != nil {
		return nil, fmt.Errorf("p.db.Query: %w", err)
	}

	posts := []*model.Post{}
	for rows.Next() {
		var strCreatedAt, strUpdatedAt string
		post := &model.Post{}
		err = rows.Scan(&post.Id, &post.Title, &post.Content, &post.UserId, &strCreatedAt, &strUpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}

		timeCreatedAt, err := time.ParseInLocation(time.RFC3339, strCreatedAt, time.Local)
		if err != nil {
			return nil, fmt.Errorf("time.Parse: %w", err)
		}
		post.CreatedAt = timeCreatedAt

		timeUpdatedAt, err := time.ParseInLocation(time.RFC3339, strUpdatedAt, time.Local)
		if err != nil {
			return nil, fmt.Errorf("time.Parse: %w", err)
		}
		post.UpdatedAt = timeUpdatedAt

		posts = append(posts, post)
	}

	return posts, nil
}

func (p *PostRepo) GetByIDs(ids []int64) ([]*model.Post, error) {
	strIDs := strings.Trim(strings.Replace(fmt.Sprint(ids), " ", ",", -1), "[]")
	preQuery := fmt.Sprintf(`SELECT id, title, content, user_id, created_at, updated_at FROM posts
WHERE id IN (%v)`, strIDs)

	rows, err := p.db.Query(preQuery)
	if err != nil {
		return nil, fmt.Errorf("db.Query: %w", err)
	}
	posts := []*model.Post{}
	for rows.Next() {
		var strCreatedAt, strUpdatedAt string
		post := &model.Post{}
		err = rows.Scan(&post.Id, &post.Title, &post.Content, &post.UserId, &strCreatedAt, &strUpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}

		timeCreatedAt, err := time.ParseInLocation(time.RFC3339, strCreatedAt, time.Local)
		if err != nil {
			return nil, fmt.Errorf("time.Parse: %w", err)
		}
		post.CreatedAt = timeCreatedAt

		timeUpdatedAt, err := time.ParseInLocation(time.RFC3339, strUpdatedAt, time.Local)
		if err != nil {
			return nil, fmt.Errorf("time.Parse: %w", err)
		}
		post.UpdatedAt = timeUpdatedAt

		posts = append(posts, post)
	}
	return posts, nil
}
func (p *PostRepo) GetByID(id int64) (*model.Post, error) {
	row := p.db.QueryRow(`
	SELECT id, title, content, user_id, created_at, updated_at FROM posts
	WHERE id = ?`, id)

	post := &model.Post{}
	var strCreatedAt, strUpdatedAt string
	err := row.Scan(&post.Id, &post.Title, &post.Content, &post.UserId, &strCreatedAt, &strUpdatedAt)
	switch {
	case err == nil:
	case strings.HasPrefix(err.Error(), "sql: no rows in result set"):
		return nil, model.ErrPostNotFound
	case err != nil:
		return nil, fmt.Errorf("row.Scan: %w", err)
	}

	timeCreatedAt, err := time.ParseInLocation(time.RFC3339, strCreatedAt, time.Local)
	if err != nil {
		return nil, fmt.Errorf("time.Parse: %w", err)
	}
	post.CreatedAt = timeCreatedAt

	timeUpdatedAt, err := time.ParseInLocation(time.RFC3339, strUpdatedAt, time.Local)
	if err != nil {
		return nil, fmt.Errorf("time.Parse: %w", err)
	}
	post.UpdatedAt = timeUpdatedAt
	return post, nil
}

func (p *PostRepo) GetAll(offset, limit int64) ([]*model.Post, error) {
	if limit == 0 {
		limit = -1
	}

	rows, err := p.db.Query(`
SELECT id, title, content, user_id, created_at, updated_at FROM posts
LIMIT ? OFFSET ? 
	`, limit, offset)

	if err != nil {
		return nil, fmt.Errorf("p.db.Query: %w", err)
	}

	posts := []*model.Post{}
	for rows.Next() {
		var strCreatedAt, strUpdatedAt string
		post := &model.Post{}
		err = rows.Scan(&post.Id, &post.Title, &post.Content, &post.UserId, &strCreatedAt, &strUpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}

		timeCreatedAt, err := time.ParseInLocation(time.RFC3339, strCreatedAt, time.Local)
		if err != nil {
			return nil, fmt.Errorf("time.Parse: %w", err)
		}
		post.CreatedAt = timeCreatedAt

		timeUpdatedAt, err := time.ParseInLocation(time.RFC3339, strUpdatedAt, time.Local)
		if err != nil {
			return nil, fmt.Errorf("time.Parse: %w", err)
		}
		post.UpdatedAt = timeUpdatedAt

		posts = append(posts, post)
	}

	return posts, nil
}

func (p *PostRepo) Update(post *model.Post) error {
	strUpdatedAt := post.UpdatedAt.Format(time.RFC3339)
	_, err := p.db.Exec(`UPDATE posts
SET title = ?, content = ?, updated_at = ?
WHERE id = ?`, post.Title, post.Content, strUpdatedAt, post.Id)
	switch {
	case err == nil:
	case err != nil:
		return fmt.Errorf("p.db.Exec: %w", err)
	}
	return nil
}

func (p *PostRepo) Create(post *model.Post) (int64, error) {
	strCreatedAt := post.CreatedAt.Format(time.RFC3339)
	row := p.db.QueryRow(`
INSERT INTO posts (title, content, user_id, created_at, updated_at) VALUES
(?, ?, ?, ?, ?) RETURNING id`, post.Title, post.Content, post.UserId, strCreatedAt, strCreatedAt)

	err := row.Scan(&post.Id)
	switch {
	case err == nil:
		return post.Id, nil
	case strings.HasPrefix(err.Error(), "constraint failed: CHECK constraint failed"):
		// Create Error
		switch {
		case strings.Contains(err.Error(), "title"):
			switch {
			case strings.Contains(err.Error(), "LENGTH"):
				return -1, model.ErrInvalidTitleLength
			}
		}
	}
	return -1, fmt.Errorf("row.Scan: %w", err)
}
func (p *PostRepo) DeleteByID(id int64) error {
	_, err := p.db.Exec("DELETE FROM posts WHERE id = ?", id)
	switch {
	case err == nil:
	case err != nil:
		return fmt.Errorf("p.db.Exec: %w", err)
	}
	return nil
}
