package service

import (
	"errors"
	"fmt"
	"forum/internal/model"
	"log"
)

func (s *Service) FillPost(post *model.Post, sesUserId int64) error {
	var err error

	post.Categories, err = s.Category.GetByPostID(post.Id)
	switch {
	case err != nil:
		log.Printf("FillPost: PostCategory.GetByPostID(postId: %v): %v", post.Id, err)
	}

	post.User, err = s.User.GetByID(post.UserId)
	switch {
	case err != nil:
		log.Printf("FillPost: User.GetByID(userId: %v): %v", post.UserId, err)
	}

	like, dislike, err := s.PostReaction.GetByPostID(post.Id)
	switch {
	case err != nil:
		log.Printf("FillPost: PostReaction.GetByPostID(id: %v): %v", post.Id, err)
	}
	post.Like = like
	post.Dislike = dislike

	if sesUserId == 0 {
		return nil
	}

	rUser, err := s.PostReaction.GetPostUserReaction(sesUserId, post.Id)
	switch {
	case err == nil:
		post.UserReaction = rUser.Reaction
	case errors.Is(err, model.ErrPostReactNotFound):
	case err != nil:
		log.Printf("FillPost: PostVote.GetPostUserReaction(userId: %v, postId: %v): %v", sesUserId, post.Id, err)
	}
	return nil
}

func (s *Service) FillPosts(posts []*model.Post, sesUserId int64) error {
	for _, post := range posts {
		err := s.FillPost(post, sesUserId)
		if err != nil {
			return fmt.Errorf("FillPosts: %w", err)
		}
	}
	return nil
}
