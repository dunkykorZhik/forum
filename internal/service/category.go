package service

import (
	"errors"
	"fmt"
	"forum/internal/model"
	"forum/internal/repository"
	"time"
)

type CategoryService struct {
	repo repository.Category
}

func NewPostCategoryService(repo repository.Category) *CategoryService {
	return &CategoryService{repo}
}

func (c *CategoryService) GetPostIDsContainedCatIDs(ids []int64, offset, limit int64) ([]int64, error) {
	cats, err := c.repo.GetPostIDsContainedCatIDs(ids, offset, limit)
	switch {
	case err == nil:
	case err != nil:
		return nil, fmt.Errorf("GetPostIDsContainedCatIDs: %w", err)
	}
	return cats, nil
}
func (c *CategoryService) GetByPostID(postId int64) ([]*model.Category, error) {
	categories, err := c.repo.GetByPostID(postId)
	switch {
	case err == nil:
		return categories, nil
	}
	return nil, fmt.Errorf("c.repo.GetByPostID: %w", err)
}
func (c *CategoryService) GetByNames(names []string) ([]*model.Category, error) {
	cats, err := c.repo.GetByNames(names)
	switch {
	case err == nil:
	case err != nil:
		return nil, fmt.Errorf("GetByNames: %w", err)
	}
	return cats, nil
}
func (p *CategoryService) GetAll(offset, limit int64) ([]*model.Category, error) {
	categories, err := p.repo.GetAll(offset, limit)
	switch {
	case err == nil:
	case err != nil:
		return nil, fmt.Errorf("p.repo.GetAll: %w", err)
	}
	return categories, nil
}
func (c *CategoryService) DeleteByPostID(postId int64) error {
	err := c.repo.DeleteByPostID(postId)
	switch {
	case err == nil:
	case err != nil:
		return fmt.Errorf("c.repo.DeleteByPostID: %w", err)
	}
	return nil
}

func (c *CategoryService) AddToPostByNames(names []string, postId int64) error {
	if len(names) == 0 {
		return nil
	} else if len(names) > 5 {
		return model.ErrCategoryLimitForPost
	}

	var ids []int64 = make([]int64, len(names))
	for i, name := range names {
		cat := &model.Category{Name: name, CreatedAt: time.Now()}
		id, err := c.repo.Create(cat)
		switch {
		case err == nil:
			ids[i] = id
			continue
		case errors.Is(err, model.ErrCategoryExistName):
		case errors.Is(err, model.ErrCheckLengthName):
			return model.ErrCheckLengthName
		default:
			return fmt.Errorf("c.repo.Create: %w", err)
		}

		cat, err = c.repo.GetByName(name)
		switch {
		case err == nil:
			ids[i] = cat.Id
			continue
		default:
			return fmt.Errorf("c.repo.GetByName: %w", err)
		}
	}

	for _, id := range ids {
		_, err := c.repo.AddToPost(id, postId)
		switch {
		case err == nil:
			continue
		default:
			return fmt.Errorf("c.repo.AddToPost: %w", err)
		}
	}
	return nil
}
