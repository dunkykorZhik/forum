package service

import (
	"errors"
	"fmt"
	"forum/internal/model"
	"forum/internal/repository"
	"time"
)

type PostCommentReactService struct {
	repo repository.PostCommentReact
}

func NewPostCommentReactService(postCommentVote repository.PostCommentReact) *PostCommentReactService {
	return &PostCommentReactService{postCommentVote}
}

func (c *PostCommentReactService) GetByCommentID(commentId int64) (int64, int64, error) {
	up, down, err := c.repo.GetByCommentID(commentId)
	switch {
	case err == nil:
	case err != nil:
		return 0, 0, fmt.Errorf("c.repo.GetByCommentID: %w", err)
	}
	return up, down, nil
}
func (c *PostCommentReactService) GetCommentUserReaction(userId, commentId int64) (*model.PostCommentReaction, error) {
	pVote, err := c.repo.GetCommentUserReaction(userId, commentId)
	switch {
	case err == nil:
	case errors.Is(err, model.ErrCommReactionNotFound):
		return nil, model.ErrCommReactionNotFound
	case err != nil:
		return nil, fmt.Errorf("repo.GetCommentUserReaction: %w", err)
	}
	return pVote, nil
}
func (p *PostCommentReactService) Record(react *model.PostCommentReaction) error {
	if react.Reaction < -1 || 1 < react.Reaction {
		return model.ErrInvalidCommReaction
	}

	react.CreatedAt = time.Now()
	_, err := p.repo.Create(react)
	switch {
	case err == nil:
		return nil
	case errors.Is(err, model.ErrCommReactionExists):
	case errors.Is(err, model.ErrCommReactionNotFound):
		return model.ErrCommReactionNotFound
	case err != nil:
		return fmt.Errorf("p.repo.Create: %w", err)
	}

	react.UpdatedAt = time.Now()
	err = p.repo.Update(react)
	switch {
	case err == nil:
	case errors.Is(err, model.ErrCommReactionNotFound):
		return model.ErrCommReactionNotFound
	case err != nil:
		return fmt.Errorf("p.repo.Update: %w", err)
	}
	return nil
}
