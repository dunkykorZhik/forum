package service

import (
	"errors"
	"fmt"
	"forum/internal/model"
	"forum/internal/repository"
	"time"
)

type PostCommentService struct {
	repo repository.PostComment
}

func NewPostCommentService(repo repository.PostComment) *PostCommentService {
	return &PostCommentService{repo}
}

func (c *PostCommentService) GetByID(id int64) (*model.PostComment, error) {
	comment, err := c.repo.GetByID(id)
	switch {
	case err == nil:
	case errors.Is(err, model.ErrCommNotFound): // check if post/comm not founc
		return nil, model.ErrCommNotFound
	case err != nil:
		return nil, fmt.Errorf("c.repo.GetByID: %w", err)
	}
	return comment, nil
}

func (c *PostCommentService) GetAllByPostID(postId, offset, limit int64) ([]*model.PostComment, error) {
	comments, err := c.repo.GetAllByPostID(postId, offset, limit)
	switch {
	case err == nil:
	case err != nil:
		return nil, fmt.Errorf("c.repo.GetAllByPostID: %w", err)
	}
	return comments, nil
}

func (c *PostCommentService) Create(comment *model.PostComment) (int64, error) {
	comment.Prepare()

	if comment.ValidateContent() != nil {
		return -1, model.ErrInvalidContentLength
	}

	comment.CreatedAt = time.Now()
	commentId, err := c.repo.Create(comment)
	switch {
	case err == nil:
	case err != nil:
		return -1, fmt.Errorf("c.repo.Create: %w", err)
	}
	return commentId, nil
}
func (c *PostCommentService) DeleteByID(id int64) error {
	err := c.repo.DeleteByID(id)
	switch {
	case err == nil:
	case err != nil:
		return fmt.Errorf("c.repo.DeleteByID: %w", err)
	}
	return nil
}
