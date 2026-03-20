package server

import (
	"errors"
	"forum/internal/model"
	"log"
	"net/http"
	"strconv"
)

func (m *mainHandler) IndexHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%-30v | %-7v | %-30v \n", r.URL, r.Method, "IndexHandler")

	// Allowed Methods
	switch r.Method {
	case http.MethodGet:
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	// Logic
	removeRedirectCookie(w, r)
	switch r.Method {
	case http.MethodGet:
		posts, err := m.service.Post.GetAll(0, -1)
		if err != nil {
			log.Printf("IndexHandler: Post.GetAll: %v\n", err)
		}

		cookie := getSessionCookie(w, r)
		if cookie == nil {
			err = m.service.FillPosts(posts, 0)
			if err != nil {
				log.Printf("IndexHandler: FillPosts: %v\n", err)
			}
			pg := &Page{Posts: posts}
			m.tmpl.executeTemplate(w, pg, "home.html")
			return
		}

		session, err := m.service.Session.GetByUuid(cookie.Value)
		switch {
		case err == nil:
		case errors.Is(err, model.ErrExpired) || errors.Is(err, model.ErrSessNotFound):
			removeSessionCookie(w, r)
			err = m.service.FillPosts(posts, 0)
			if err != nil {
				log.Printf("IndexHandler: FillPosts: %v\n", err)
			}
			pg := &Page{Posts: posts}
			m.tmpl.executeTemplate(w, pg, "home.html")
			return
		case err != nil:
			log.Printf("IndexHandler: m.service.Session.GetByUuid: %v\n", err)
			http.Error(w, "something wrong, maybe try again later", http.StatusInternalServerError)
			return
		}

		user, err := m.service.User.GetByID(session.UserId)
		switch {
		case err == nil:
		case err != nil:
			log.Printf("IndexHandler: m.service.Session.GetByUuid: %v\n", err)
			http.Error(w, "something wrong, maybe try again later", http.StatusInternalServerError)
			return
		}

		err = m.service.FillPosts(posts, user.Id)
		if err != nil {
			log.Printf("IndexHandler: FillPosts: %v\n", err)
		}
		pg := &Page{Posts: posts, User: user}
		m.tmpl.executeTemplate(w, pg, "home.html")
		return
	}
}

func (m *mainHandler) PostViewHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%-30v | %-7v | %-30v \n", r.URL, r.Method, "PostViewHandler")

	// Allowed Methods
	switch r.Method {
	case http.MethodGet:
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	var user *model.User
	cookie := getSessionCookie(w, r)
	switch cookie {
	case nil:
		user = nil
	default:
		session, err := m.service.Session.GetByUuid(cookie.Value)
		switch {
		case err == nil:
			user, _ = m.service.User.GetByID(session.UserId)
		case errors.Is(err, model.ErrExpired) || errors.Is(err, model.ErrNotFound):
			removeSessionCookie(w, r)
		default:
			log.Printf("PostViewHandler: m.service.Session.GetByUuid: %v\n", err)
			http.Error(w, "something wrong, maybe try again later", http.StatusInternalServerError)
			return
		}
	}

	switch r.Method {
	case http.MethodGet:
		strPostId := r.URL.Query().Get("id")
		postId, err := strconv.ParseInt(strPostId, 10, 64)
		if err != nil || postId < 1 {
			http.Error(w, "Invalid query id", http.StatusBadRequest)
			return
		}

		post, err := m.service.Post.GetByID(postId)
		switch {
		case err == nil:
		case errors.Is(err, model.ErrPostNotFound):
			// TODO: error Page
			http.Error(w, "Post Not Found", http.StatusNotFound)
			return
		}

		var userId int64
		if user != nil {
			userId = user.Id
		}
		m.service.FillPost(post, userId)

		post.Comments, err = m.service.PostComment.GetAllByPostID(post.Id, 0, -1)
		switch {
		case err == nil:
			for _, comment := range post.Comments {
				comment.User, err = m.service.User.GetByID(comment.UserId)
				if err != nil {
					log.Printf("PostViewHandler: m.service.User.GetByID: %w", err)
				}
			}
		case err != nil:
			log.Printf("PostViewHandler: service.PostComment.GetAllByPostID: %v\n", err)
		}

		for _, comment := range post.Comments {
			comment.Like, comment.Dislike, err = m.service.PostCommentReaction.GetByCommentID(comment.Id)
			if err != nil {
				log.Printf("PostViewHandler: service.PostCommentReaction.GetByCommentID(commentId: %v): %v\n", comment.Id, err)
			}
			vt, err := m.service.PostCommentReaction.GetCommentUserReaction(userId, comment.Id)
			switch {
			case err == nil:
				comment.UserReaction = vt.Reaction
			case errors.Is(err, model.ErrCommReactionNotFound):
			case err != nil:
				log.Printf("PostViewHandler: service.PostCommentReaction.GetCommentUserVote: %v\n", err)
			}
		}

		pg := &Page{User: user, Post: post}
		m.tmpl.executeTemplate(w, pg, "post-view.html")
		return
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}
