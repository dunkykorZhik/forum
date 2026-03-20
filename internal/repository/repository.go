package repository

import (
	"database/sql"
	"forum/internal/model"
)

type Repository struct {
	User
	Session
	Post
	PostComment
	Category
	PostReaction
	PostCommentReact
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		User:             NewUserRepo(db),
		Session:          NewSessionRepo(db),
		Post:             NewPostRepo(db),
		PostComment:      NewPostCommentRepo(db),
		Category:         NewCategoryRepo(db),
		PostReaction:     NewPostReactionRepo(db),
		PostCommentReact: NewPostCommentReactionRepo(db),
	}
}

type User interface {
	Create(user *model.User) (int64, error)
	//Update(user *model.User) error
	//DeleteByID(id int64) error

	GetByID(id int64) (*model.User, error)
	GetByUsername(username string) (*model.User, error)
	GetByEmail(email string) (*model.User, error)
	// GetAll(from, offset int64) error
}
type Session interface {
	Create(session *model.Session) (int64, error)
	//Delete(id int64) error
	GetByUuid(uuid string) (*model.Session, error)
	UpdateByUserId(userId int64, session *model.Session) error
}

type PostReaction interface {
	Create(vote *model.PostReaction) (int64, error)
	Update(vote *model.PostReaction) error
	GetByPostID(postId int64) (int64, int64, error)
	GetPostUserReaction(userId, postId int64) (*model.PostReaction, error)
	GetAllUserReactedPostIDs(userId int64, reaction int8, limit, offset int64) ([]int64, error)
	//DeleteByID(id int64) error
}

type PostComment interface {
	Create(comment *model.PostComment) (int64, error)
	// Update(comment *model.PostComment) error
	GetAllByPostID(postId, offset, limit int64) ([]*model.PostComment, error)
	GetByID(id int64) (*model.PostComment, error)
	DeleteByID(id int64) error
}

type PostCommentReact interface {
	Create(reaction *model.PostCommentReaction) (int64, error)
	Update(reaction *model.PostCommentReaction) error
	GetByCommentID(commentId int64) (int64, int64, error)
	GetCommentUserReaction(userId, commentId int64) (*model.PostCommentReaction, error)
	//DeleteByID(id int64) error
}

type Post interface {
	Create(post *model.Post) (int64, error)
	Update(post *model.Post) error
	GetAll(offset, limit int64) ([]*model.Post, error)
	GetByID(id int64) (*model.Post, error)
	GetByIDs(ids []int64) ([]*model.Post, error)
	GetByUserID(userId, offset, limit int64) ([]*model.Post, error)
	DeleteByID(id int64) error
}
type Category interface {
	Create(category *model.Category) (int64, error)
	AddToPost(categoryId, postId int64) (int64, error)
	//Update(category *model.Category) error
	//	GetByID(id int64) (*model.Category, error)
	GetByName(name string) (*model.Category, error)
	GetByNames(names []string) ([]*model.Category, error)
	GetByPostID(postId int64) ([]*model.Category, error)
	GetPostIDsContainedCatIDs(ids []int64, offset, limit int64) ([]int64, error)
	GetAll(offset, limit int64) ([]*model.Category, error)
	DeleteByPostID(postId int64) error
	//	DeleteByID(id int64) error
}
