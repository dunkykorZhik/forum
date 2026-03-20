package model

import (
	"errors"
	"time"
)

var (
	ErrInvalidUsername     = errors.New("invalid nickname")
	ErrInvalidEmail        = errors.New("invalid email")
	ErrExistUsername       = errors.New("user with this userame exists")
	ErrExistEmail          = errors.New("user with this email exists")
	ErrWrongLengthUsername = errors.New("user nickname length is wrong")
	ErrWrongLengthEmail    = errors.New("user email length is wrong")
	ErrNotFound            = errors.New("user not found")

	ErrSessNotFound  = errors.New("session not found")
	ErrExpired       = errors.New("session time expired")
	ErrSessionExists = errors.New("session already exists")

	ErrPostNotFound         = errors.New("post not found")
	ErrPostReactNotFound    = errors.New("post reaction not found")
	ErrInvalidContentLength = errors.New("invalid content length")
	ErrInvalidTitleLength   = errors.New("invalid title length")
	ErrInvalidPostReaction  = errors.New("invalid post reaction")
	ErrPostReactionExists   = errors.New("reaction already exists")

	ErrCommNotFound         = errors.New("comment not found")
	ErrInvalidCommReaction  = errors.New("invalid comment reaction")
	ErrCommReactionNotFound = errors.New("comment reaction not found")
	ErrCommReactionExists   = errors.New("reaction already exists")

	ErrCheckLengthName      = errors.New("category namelength is wrong")
	ErrCategoryExistName    = errors.New("category already exists")
	ErrCategoryNotFound     = errors.New("category not found")
	ErrCategoryLimitForPost = errors.New("number of categories are limited")
)

type User struct {
	Id        int64
	Username  string
	Email     string
	Password  string
	CreatedAt time.Time
}

type Session struct {
	Id        int64
	Uuid      string
	ExpiredAt time.Time
	UserId    int64
}

type Post struct {
	Id        int64
	Title     string
	Content   string
	UserId    int64
	CreatedAt time.Time
	UpdatedAt time.Time

	User         *User
	UserReaction int8  // -1 0 1
	Like         int64 // Like
	Dislike      int64 // Dislike
	Categories   []*Category
	Comments     []*PostComment
}

type PostComment struct {
	Id        int64
	Content   string
	PostId    int64
	UserId    int64
	CreatedAt time.Time

	User         *User
	UserReaction int8  // -1 0 1
	Like         int64 // Like
	Dislike      int64 // Dislike
}

type Category struct {
	Id        int64
	Name      string
	CreatedAt time.Time
}

type PostReaction struct {
	Id        int64
	PostId    int64
	UserId    int64
	Reaction  int8
	CreatedAt time.Time
	UpdatedAt time.Time
}

type PostCommentReaction struct {
	Id        int64
	CommentId int64
	UserId    int64
	Reaction  int8
	CreatedAt time.Time
	UpdatedAt time.Time
}
