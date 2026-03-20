package service

import (
	"errors"
	"fmt"
	"forum/internal/model"
	"forum/internal/repository"
	"time"
)

type PostService struct {
	repo repository.Post
}

func NewPostService(repo repository.Post) *PostService {
	return &PostService{repo}
}

func (p *PostService) GetByID(id int64) (*model.Post, error) {
	post, err := p.repo.GetByID(id)
	switch {
	case err == nil:
		return post, nil
	case errors.Is(err, model.ErrPostNotFound):
		return nil, model.ErrPostNotFound
	}
	return nil, fmt.Errorf("p.repo.GetByID: %w", err)
}

func (p *PostService) GetByIDs(ids []int64) ([]*model.Post, error) {
	posts, err := p.repo.GetByIDs(ids)
	switch {
	case err == nil:
	case err != nil:
		return nil, fmt.Errorf("GetByIDs: %w", err)
	}
	return posts, nil
}
func (p *PostService) GetByUserID(userId, offset, limit int64) ([]*model.Post, error) {
	posts, err := p.repo.GetByUserID(userId, offset, limit)
	switch {
	case err == nil:
	case err != nil:
		return nil, fmt.Errorf("p.repo.GetByUserID: %w", err)
	}
	return posts, nil
}

func (p *PostService) Update(post *model.Post) error {
	post.Prepare()

	if post.ValidateTitle() != nil {
		return model.ErrInvalidTitleLength
	} else if post.ValidateContent() != nil {
		return model.ErrInvalidContentLength
	}

	post.UpdatedAt = time.Now()
	err := p.repo.Update(post)
	switch {
	case err == nil:
	case err != nil:
		return fmt.Errorf("p.repo.Update: %w", err)
	}
	return nil
}

func (p *PostService) Create(post *model.Post) (int64, error) {
	post.Prepare()

	if post.ValidateTitle() != nil {
		return -1, model.ErrInvalidTitleLength
	} else if post.ValidateContent() != nil {
		return -1, model.ErrInvalidContentLength
	}

	post.CreatedAt = time.Now()
	post.UpdatedAt = post.CreatedAt

	postId, err := p.repo.Create(post)
	switch {
	case err == nil:
	case errors.Is(err, model.ErrInvalidTitleLength):
		return -1, model.ErrInvalidTitleLength
	case err != nil:
		return -1, fmt.Errorf("p.repo.Create: %w", err)
	}
	return postId, nil

}

func (q *PostService) DeleteByID(id int64) error {
	err := q.repo.DeleteByID(id)
	switch {
	case err == nil:
	case err != nil:
		return fmt.Errorf("q.repo.DeleteByID: %w", err)
	}
	return nil
}

func (p *PostService) GetAll(offset, limit int64) ([]*model.Post, error) {
	posts, err := p.repo.GetAll(offset, limit)
	switch {
	case err == nil:
	case err != nil:
		return nil, fmt.Errorf("p.repo.GetAll: %w", err)
	}
	return posts, nil
}
