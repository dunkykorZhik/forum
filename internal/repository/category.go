package repository

import (
	"database/sql"
	"fmt"
	"forum/internal/model"
	"strings"
	"time"
)

type CategoryRepo struct {
	db *sql.DB
}

func NewCategoryRepo(db *sql.DB) *CategoryRepo {
	return &CategoryRepo{db}
}
func (p *CategoryRepo) GetAll(offset, limit int64) ([]*model.Category, error) {
	if limit == 0 {
		limit = -1
	}

	rows, err := p.db.Query(`
SELECT id, name, created_at FROM categories
LIMIT ? OFFSET ? 
	`, limit, offset)

	if err != nil {
		return nil, fmt.Errorf("p.db.Query: %w", err)
	}

	categories := []*model.Category{}
	for rows.Next() {
		var strCreatedAt string
		category := &model.Category{}
		err = rows.Scan(&category.Id, &category.Name, &strCreatedAt)
		if err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}

		timeCreatedAt, err := time.ParseInLocation(time.RFC3339, strCreatedAt, time.Local)
		if err != nil {
			return nil, fmt.Errorf("time.Parse: %w", err)
		}
		category.CreatedAt = timeCreatedAt

		categories = append(categories, category)
	}

	return categories, nil
}
func (c *CategoryRepo) GetByName(name string) (*model.Category, error) {
	row := c.db.QueryRow(`
SELECT id, name, created_at FROM categories
WHERE name = ?`, name)

	category := &model.Category{}
	var strCreatedAt string
	err := row.Scan(&category.Id, &category.Name, &strCreatedAt)
	switch {
	case err == nil:
		timeCreatedAt, err := time.ParseInLocation(time.RFC3339, strCreatedAt, time.Local)
		if err != nil {
			return nil, fmt.Errorf("time.Parse: %w", err)
		}
		category.CreatedAt = timeCreatedAt
		return category, nil
	case strings.HasPrefix(err.Error(), "sql: no rows in result set"):
		return nil, model.ErrCategoryNotFound
	}
	return nil, fmt.Errorf("row.Scan: %w", err)
}

func (c *CategoryRepo) GetByNames(names []string) ([]*model.Category, error) {
	if len(names) == 0 {
		return nil, nil
	}

	iNames := make([]interface{}, len(names))
	for i, v := range names {
		iNames[i] = v
	}

	strQuery := fmt.Sprintf(`SELECT id, name, created_at FROM categories
WHERE name IN (%v)`, `?`+strings.Repeat(",?", len(iNames)-1))
	rows, err := c.db.Query(strQuery, iNames...)
	if err != nil {
		return nil, fmt.Errorf("db.Query: %w", err)
	}

	cats := []*model.Category{}
	for rows.Next() {
		var strCreatedAt string
		cat := &model.Category{}
		err = rows.Scan(&cat.Id, &cat.Name, &strCreatedAt)
		if err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}

		timeCreatedAt, err := time.ParseInLocation(time.RFC3339, strCreatedAt, time.Local)
		if err != nil {
			return nil, fmt.Errorf("time.Parse: %w", err)
		}
		cat.CreatedAt = timeCreatedAt

		cats = append(cats, cat)
	}
	return cats, nil
}

func (c *CategoryRepo) GetByPostID(postId int64) ([]*model.Category, error) {
	rows, err := c.db.Query(`
SELECT c.id, c.name, c.created_at FROM posts_categories pc
JOIN categories c ON pc.category_id = c.id
WHERE pc.post_id = ?`, postId)
	if err != nil {
		return nil, fmt.Errorf("c.db.Query: %w", err)
	}

	categories := []*model.Category{}
	for rows.Next() {
		var strCreatedAt string
		category := &model.Category{}
		err = rows.Scan(&category.Id, &category.Name, &strCreatedAt)
		if err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}

		timeCreatedAt, err := time.ParseInLocation(time.RFC3339, strCreatedAt, time.Local)
		if err != nil {
			return nil, fmt.Errorf("time.Parse: %w", err)
		}
		category.CreatedAt = timeCreatedAt
		categories = append(categories, category)
	}
	return categories, nil
}

func (c *CategoryRepo) GetPostIDsContainedCatIDs(ids []int64, offset, limit int64) ([]int64, error) {
	strIDs := strings.Trim(strings.Replace(fmt.Sprint(ids), " ", ",", -1), "[]")
	preQuery := fmt.Sprintf(`SELECT post_id, COUNT(category_id) as cat from posts_categories
WHERE category_id IN (%s)
GROUP BY post_id
HAVING cat >= %d
LIMIT ? OFFSET ?`, strIDs, len(ids))

	rows, err := c.db.Query(preQuery, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("db.Query: %w", err)
	}

	postIDs := []int64{}
	for rows.Next() {
		var postId, a int64
		err = rows.Scan(&postId, &a)
		if err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}

		postIDs = append(postIDs, postId)
	}
	return postIDs, nil
}
func (c *CategoryRepo) AddToPost(categoryId, postId int64) (int64, error) {
	row := c.db.QueryRow(`
	INSERT INTO posts_categories (post_id, category_id) VALUES
	(?, ?) RETURNING id`, postId, categoryId)

	id := int64(-1)
	err := row.Scan(&id)
	switch {
	case err == nil:
		return id, nil
	}
	return -1, fmt.Errorf("row.Scan: %w", err)
}

func (c *CategoryRepo) Create(category *model.Category) (int64, error) {
	strCreatedAt := category.CreatedAt.Format(time.RFC3339)
	row := c.db.QueryRow(`
	INSERT INTO categories (name, created_at) VALUES
	(?, ?) RETURNING id`, category.Name, strCreatedAt)
	err := row.Scan(&category.Id)

	switch {
	case err == nil:
		return category.Id, nil
	case strings.HasPrefix(err.Error(), "constraint failed: UNIQUE constraint failed"):
		return -1, model.ErrCategoryExistName
	case strings.HasPrefix(err.Error(), "constraint failed: CHECK constraint failed"):
		return -1, model.ErrCheckLengthName
	}
	return -1, fmt.Errorf("row.Scan: %w", err)
}
func (c *CategoryRepo) DeleteByPostID(postId int64) error {
	_, err := c.db.Exec(`DELETE FROM posts_categories
WHERE post_id = ?`, postId)
	switch {
	case err == nil:
	case err != nil:
		return fmt.Errorf("c.db.Exec: %w", err)
	}
	return nil
}
