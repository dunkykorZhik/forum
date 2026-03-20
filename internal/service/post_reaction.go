package service

import (
	"errors"
	"fmt"
	"forum/internal/model"
	"forum/internal/repository"
	"time"
)

type PostReactService struct {
	repo repository.PostReaction
}

func NewPostReactService(postVote repository.PostReaction) *PostReactService {
	return &PostReactService{postVote}
}

func (p *PostReactService) GetPostUserReaction(userId, postId int64) (*model.PostReaction, error) {
	pReaction, err := p.repo.GetPostUserReaction(userId, postId)
	switch {
	case err == nil:
	case errors.Is(err, model.ErrPostReactNotFound):
		return nil, model.ErrPostReactNotFound
	case err != nil:
		return nil, fmt.Errorf("p.repo.GetPostUserReaction: %w", err)
	}
	return pReaction, nil
}
func (p *PostReactService) GetByPostID(postId int64) (int64, int64, error) {
	up, down, err := p.repo.GetByPostID(postId)
	switch {
	case err == nil:
	case err != nil:
		return 0, 0, fmt.Errorf("p.repo.GetByPostID: %w", err)
	}
	return up, down, nil
}
func (p *PostReactService) GetAllUserReactedPostIDs(userId int64, react int8, limit, offset int64) ([]int64, error) {
	postIDs, err := p.repo.GetAllUserReactedPostIDs(userId, react, limit, offset)
	switch {
	case err == nil:
	case err != nil:
		return nil, fmt.Errorf("GetAllUserVotedPostIDs: %w", err)
	}
	return postIDs, nil
}
func (p *PostReactService) Record(react *model.PostReaction) error {
	if react.Reaction < -1 || 1 < react.Reaction {
		return model.ErrInvalidPostReaction
	}

	react.CreatedAt = time.Now()
	_, err := p.repo.Create(react)
	switch {
	case err == nil:
		return nil
	case errors.Is(err, model.ErrPostReactionExists):
	case errors.Is(err, model.ErrPostReactNotFound):
		return model.ErrPostReactNotFound
	case err != nil:
		return fmt.Errorf("p.repo.Create: %w", err)
	}

	react.UpdatedAt = time.Now()
	err = p.repo.Update(react)
	switch {
	case err == nil:
	case errors.Is(err, model.ErrPostReactNotFound):
		return model.ErrPostReactNotFound
	case err != nil:
		return fmt.Errorf("p.repo.Update: %w", err)
	}
	return nil
}
