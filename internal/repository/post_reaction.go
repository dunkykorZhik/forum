package repository

import (
	"database/sql"
	"fmt"
	"forum/internal/model"
	"strings"
	"time"
)

type PostReactionRepo struct {
	db *sql.DB
}

func NewPostReactionRepo(db *sql.DB) *PostReactionRepo {
	return &PostReactionRepo{db}
}

func (p *PostReactionRepo) GetByPostID(postId int64) (int64, int64, error) {
	row := p.db.QueryRow(`
SELECT 
	(SELECT COUNT(vote) FROM posts_votes WHERE post_id = ? AND vote == 1) as up,
    (SELECT COUNT(vote) FROM posts_votes WHERE post_id = ? AND vote == -1) as down;
	`, postId, postId)

	var up, down int64
	err := row.Scan(&up, &down)
	switch {
	case err == nil:
	case err != nil:
		return 0, 0, fmt.Errorf("row.Scan: %w", err)
	}
	return up, down, nil
}

func (p *PostReactionRepo) GetPostUserReaction(userId, postId int64) (*model.PostReaction, error) {
	row := p.db.QueryRow(`
	SELECT id, post_id, user_id, vote, created_at FROM posts_votes
	WHERE post_id = ? AND user_id = ?`, postId, userId)
	postVote := &model.PostReaction{}

	strCreatedAt := ""
	err := row.Scan(&postVote.Id, &postVote.PostId, &postVote.UserId, &postVote.Reaction, &strCreatedAt)

	switch {
	case err == nil:
	case strings.HasPrefix(err.Error(), "sql: no rows in result set"):
		return nil, model.ErrPostReactNotFound
	default:
		return nil, fmt.Errorf("row.Scan: %w", err)
	}

	timeCreatedAt, err := time.ParseInLocation(time.RFC3339, strCreatedAt, time.Local)
	if err != nil {
		return nil, fmt.Errorf("time.Parse: %w", err)
	}
	postVote.CreatedAt = timeCreatedAt
	return postVote, nil
}

func (p *PostReactionRepo) GetAllUserReactedPostIDs(userId int64, vote int8, limit, offset int64) ([]int64, error) {
	if limit == 0 {
		limit = -1
	}

	rows, err := p.db.Query(`
	SELECT post_id FROM posts_votes
	WHERE user_id = ? AND vote = ?
	LIMIT ? OFFSET ?`, userId, vote, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("p.db.Query: %w", err)
	}

	postIDs := []int64{}
	for rows.Next() {
		var postId int64
		err = rows.Scan(&postId)
		if err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}

		postIDs = append(postIDs, postId)
	}

	return postIDs, nil
}

func (p *PostReactionRepo) Create(vote *model.PostReaction) (int64, error) {
	strCreatedAt := vote.CreatedAt.Format(time.RFC3339)
	row := p.db.QueryRow(`
INSERT INTO posts_votes (vote, user_id, post_id, created_at, updated_at) VALUES
(?, ?, ?, ?, ?) RETURNING id`, vote.Reaction, vote.UserId, vote.PostId, strCreatedAt, strCreatedAt)

	err := row.Scan(&vote.Id)
	switch {
	case err == nil:
	case err != nil:
		switch {
		case strings.HasPrefix(err.Error(), "constraint failed: UNIQUE constraint failed"):
			return -1, model.ErrPostReactionExists
		case strings.HasPrefix(err.Error(), "constraint failed: FOREIGN KEY constraint failed"):
			return -1, model.ErrPostReactNotFound
		}
		return -1, fmt.Errorf("row.Scan: %w", err)
	}
	return vote.Id, nil
}

func (p *PostReactionRepo) Update(vote *model.PostReaction) error {
	strUpdatedAt := vote.UpdatedAt.Format(time.RFC3339)
	row := p.db.QueryRow(`
UPDATE posts_votes
SET vote = ?, updated_at = ?
WHERE user_id = ? AND post_id = ? 
RETURNING id`, vote.Reaction, strUpdatedAt, vote.UserId, vote.PostId)

	err := row.Scan(&vote.Id)
	switch {
	case err == nil:
	case strings.HasPrefix(err.Error(), "constraint failed: FOREIGN KEY constraint failed"):
		return model.ErrPostReactNotFound
	case err != nil:
		return fmt.Errorf("row.Scan: %w", err)
	}
	return nil
}
