package repository

import (
	"database/sql"
	"fmt"
	"forum/internal/model"
	"strings"
	"time"
)

type PostCommentReactionRepo struct {
	db *sql.DB
}

func NewPostCommentReactionRepo(db *sql.DB) *PostCommentReactionRepo {
	return &PostCommentReactionRepo{db}
}

func (c *PostCommentReactionRepo) GetByCommentID(commentId int64) (int64, int64, error) {
	row := c.db.QueryRow(`
SELECT 
	(SELECT COUNT(vote) FROM posts_comments_votes WHERE comment_id = ? AND vote == 1) as up,
    (SELECT COUNT(vote) FROM posts_comments_votes WHERE comment_id = ? AND vote == -1) as down;
	`, commentId, commentId)

	var up, down int64
	err := row.Scan(&up, &down)
	switch {
	case err == nil:
	case err != nil:
		return 0, 0, fmt.Errorf("row.Scan: %w", err)
	}
	return up, down, nil
}
func (c *PostCommentReactionRepo) GetCommentUserReaction(userId, commentId int64) (*model.PostCommentReaction, error) {
	row := c.db.QueryRow(`
	SELECT id, comment_id, user_id, vote, created_at FROM posts_comments_votes
	WHERE comment_id = ? AND user_id = ?`, commentId, userId)
	commentVote := &model.PostCommentReaction{}

	strCreatedAt := ""
	err := row.Scan(&commentVote.Id, &commentVote.CommentId, &commentVote.UserId, &commentVote.Reaction, &strCreatedAt)

	switch {
	case err == nil:
	case strings.HasPrefix(err.Error(), "sql: no rows in result set"):
		return nil, model.ErrCommReactionNotFound
	default:
		return nil, fmt.Errorf("row.Scan: %w", err)
	}

	timeCreatedAt, err := time.ParseInLocation(time.RFC3339, strCreatedAt, time.Local)
	if err != nil {
		return nil, fmt.Errorf("time.Parse: %w", err)
	}
	commentVote.CreatedAt = timeCreatedAt
	return commentVote, nil
}
func (p *PostCommentReactionRepo) Create(vote *model.PostCommentReaction) (int64, error) {
	strCreatedAt := vote.CreatedAt.Format(time.RFC3339)
	row := p.db.QueryRow(`
INSERT INTO posts_comments_votes (vote, user_id, comment_id, created_at, updated_at) VALUES
(?, ?, ?, ?, ?) RETURNING id`, vote.Reaction, vote.UserId, vote.CommentId, strCreatedAt, strCreatedAt)

	err := row.Scan(&vote.Id)
	switch {
	case err == nil:
	case err != nil:
		switch {
		case strings.HasPrefix(err.Error(), "constraint failed: UNIQUE constraint failed"):
			return -1, model.ErrCommReactionExists
		case strings.HasPrefix(err.Error(), "constraint failed: FOREIGN KEY constraint failed"):
			return -1, model.ErrCommReactionNotFound
		}
		return -1, fmt.Errorf("row.Scan: %w", err)
	}
	return vote.Id, nil
}

func (p *PostCommentReactionRepo) Update(vote *model.PostCommentReaction) error {
	strUpdatedAt := vote.UpdatedAt.Format(time.RFC3339)
	row := p.db.QueryRow(`
UPDATE posts_comments_votes
SET vote = ?, updated_at = ?
WHERE user_id = ? AND comment_id = ? 
RETURNING id`, vote.Reaction, strUpdatedAt, vote.UserId, vote.CommentId)

	err := row.Scan(&vote.Id)
	switch {
	case err == nil:
	case strings.HasPrefix(err.Error(), "constraint failed: FOREIGN KEY constraint failed"):
		return model.ErrCommReactionNotFound
	case err != nil:
		return fmt.Errorf("row.Scan: %w", err)
	}
	return nil
}
