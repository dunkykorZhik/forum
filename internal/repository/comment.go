package repository

import (
	"database/sql"
	"fmt"
	"forum/internal/model"
	"strings"
	"time"
)

type PostCommentRepo struct {
	db *sql.DB
}

func NewPostCommentRepo(db *sql.DB) *PostCommentRepo {
	return &PostCommentRepo{db}
}

func (c *PostCommentRepo) GetByID(id int64) (*model.PostComment, error) {
	row := c.db.QueryRow(`
	SELECT id, content, user_id, post_id, created_at FROM posts_comments
	WHERE id = ?`, id)

	comment := &model.PostComment{}
	var strCreatedAt string

	err := row.Scan(&comment.Id, &comment.Content, &comment.UserId, &comment.PostId, &strCreatedAt)
	switch {
	case err == nil:
	case strings.HasPrefix(err.Error(), "sql: no rows in result set"):
		return nil, model.ErrCommNotFound
	case err != nil:
		return nil, fmt.Errorf("row.Scan: %w", err)
	}

	timeCreatedAt, err := time.ParseInLocation(time.RFC3339, strCreatedAt, time.Local)
	if err != nil {
		return nil, fmt.Errorf("time.Parse: %w", err)
	}
	comment.CreatedAt = timeCreatedAt

	return comment, nil
}
func (c *PostCommentRepo) GetAllByPostID(postId, offset, limit int64) ([]*model.PostComment, error) {
	if limit == 0 {
		limit = -1
	}

	rows, err := c.db.Query(`
SELECT id, content, user_id, post_id, created_at FROM posts_comments
WHERE post_id = ?
LIMIT ? OFFSET ? 
	`, postId, limit, offset)

	if err != nil {
		return nil, fmt.Errorf("p.db.Query: %w", err)
	}

	comments := []*model.PostComment{}
	for rows.Next() {
		var strCreatedAt string
		comment := &model.PostComment{}

		err := rows.Scan(&comment.Id, &comment.Content, &comment.UserId, &comment.PostId, &strCreatedAt)
		if err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}

		timeCreatedAt, err := time.ParseInLocation(time.RFC3339, strCreatedAt, time.Local)
		if err != nil {
			return nil, fmt.Errorf("time.Parse: %w", err)
		}
		comment.CreatedAt = timeCreatedAt
		comments = append(comments, comment)
	}
	return comments, nil
}
func (c *PostCommentRepo) Create(comment *model.PostComment) (int64, error) {
	strCreatedAt := comment.CreatedAt.Format(time.RFC3339)
	row := c.db.QueryRow(`
INSERT INTO posts_comments (content, user_id, post_id, created_at) VALUES
(?, ?, ?, ?) RETURNING id`, comment.Content, comment.UserId, comment.PostId, strCreatedAt)

	err := row.Scan(&comment.Id)
	switch {
	case err == nil:
	// case strings.HasPrefix(err.Error(), "FOREIGN KEY constraint failed"):
	// 	return -1, ErrNotFound
	case err != nil:
		return -1, fmt.Errorf("row.Scan: %w", err)
	}
	return comment.Id, nil
}
func (c *PostCommentRepo) DeleteByID(id int64) error {
	_, err := c.db.Exec("DELETE FROM posts_comments WHERE id = ?", id)
	switch {
	case err == nil:
	case err != nil:
		return fmt.Errorf("p.db.Exec: %w", err)
	}
	return nil
}
