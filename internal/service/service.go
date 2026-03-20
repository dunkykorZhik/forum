package service

import (
	"forum/internal/model"
	"forum/internal/repository"
)

type Service struct {
	User
	Session
	Post
	PostComment
	Category
	PostReaction
	PostCommentReaction
	repo *repository.Repository
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		User:                NewUserService(repo.User),
		Session:             NewSessionService(repo.Session),
		Post:                NewPostService(repo.Post),
		PostComment:         NewPostCommentService(repo.PostComment),
		Category:            NewPostCategoryService(repo.Category),
		PostReaction:        NewPostReactService(repo.PostReaction),
		PostCommentReaction: NewPostCommentReactService(repo.PostCommentReact),
	}
}

type User interface {
	Create(user *model.User) (int64, error)
	//Update(user *model.User) error
	//DeleteByID(id int64) error

	GetByID(id int64) (*model.User, error)
	GetByUsernameOrEmail(field string) (*model.User, error)
	// GetAll(from, offset int64) error
}
type Session interface {
	Record(userId int64) (*model.Session, error)
	//Delete(id int64) error
	GetByUuid(uuid string) (*model.Session, error)
}
type PostReaction interface {
	Record(vote *model.PostReaction) error
	GetByPostID(postId int64) (int64, int64, error)
	GetPostUserReaction(userId, postId int64) (*model.PostReaction, error)
	GetAllUserReactedPostIDs(userId int64, vote int8, limit, offset int64) ([]int64, error)
	//DeleteByID(id int64) error
}

type PostComment interface {
	Create(comment *model.PostComment) (int64, error)
	GetAllByPostID(postId, offset, limit int64) ([]*model.PostComment, error)
	GetByID(id int64) (*model.PostComment, error)
	DeleteByID(id int64) error

	// Update(comment *model.PostComment) error
}

type PostCommentReaction interface {
	Record(vote *model.PostCommentReaction) error
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
	AddToPostByNames(names []string, postId int64) error
	GetByPostID(postId int64) ([]*model.Category, error)
	GetByNames(names []string) ([]*model.Category, error)
	GetAll(offset, limit int64) ([]*model.Category, error)
	GetPostIDsContainedCatIDs(ids []int64, offset, limit int64) ([]int64, error)
	DeleteByPostID(postId int64) error
	//DeleteFromPost(id int64) error
}
