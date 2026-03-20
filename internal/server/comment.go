package server

import (
	"errors"
	"fmt"
	"forum/internal/model"
	"log"
	"net/http"
	"strconv"
)

func (m *mainHandler) PostCommentCreateHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%-30v | %-7v | %-30v \n", r.URL, r.Method, "PostCommentCreateHandler")

	// Allowed Methods
	switch r.Method {
	case http.MethodPost:
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	iUserId := r.Context().Value("UserId")
	if iUserId == nil {
		log.Println("PostCommentCreateHandler: r.Context().Value(\"UserId\") is nil")
		pg := &Page{Error: fmt.Errorf("internal server error, maybe try again later")}
		w.WriteHeader(http.StatusInternalServerError)
		m.tmpl.executeTemplate(w, pg, "post-create.html")
		return
	}

	userId := iUserId.(int64)
	user, err := m.service.User.GetByID(userId)
	switch {
	case err == nil:
	case errors.Is(err, model.ErrNotFound):
		removeSessionCookie(w, r)
		if http.MethodGet == r.Method {
			addRedirectCookie(w, r.RequestURI)
		}
		http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
		return
	case err != nil:
		log.Printf("PostEditHandler: m.service.User.GetByID: %v\n", err)
		pg := &Page{Error: fmt.Errorf("internal server error, maybe try again later")}
		w.WriteHeader(http.StatusInternalServerError)
		m.tmpl.executeTemplate(w, pg, "alert.html") // TODO: Custom Epage
		return
	}

	switch r.Method {
	case http.MethodPost:
		r.ParseForm()

		strPostId := r.FormValue("post_id")
		postId, err := strconv.ParseInt(strPostId, 10, 64)
		if err != nil || postId < 1 {
			http.Error(w, "Invalid query id", http.StatusBadRequest)
			return
		}

		comment := &model.PostComment{
			Content: r.FormValue("content"),
			PostId:  postId,
			UserId:  user.Id,
		}
		_, err = m.service.PostComment.Create(comment)
		switch {
		case err == nil:
		case errors.Is(err, model.ErrInvalidContentLength):
			http.Error(w, "invalid length of content", http.StatusBadRequest)
			return
		default:
			log.Printf("PostCommentCreateHandler: m.service.PostComment.Create: %s", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
		return
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (m *mainHandler) PostCommentDeleteHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%-30v | %-7v | %-30v \n", r.URL, r.Method, "PostCommentDelete")

	// Allowed Methods
	switch r.Method {
	case http.MethodGet:
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	iUserId := r.Context().Value("UserId")
	if iUserId == nil {
		log.Println("PostCommentDeleteHandler: r.Context().Value(\"UserId\") is nil")
		pg := &Page{Error: fmt.Errorf("internal server error, maybe try again later")}
		w.WriteHeader(http.StatusInternalServerError)
		m.tmpl.executeTemplate(w, pg, "alert.html") // TODO: Custom Epage
		return
	}

	userId := iUserId.(int64)

	switch r.Method {
	case http.MethodGet:
		strPostCommentId := r.URL.Query().Get("id")
		postCommentId, err := strconv.ParseInt(strPostCommentId, 10, 64)
		if err != nil || postCommentId < 1 {
			http.Error(w, "Invalid query id", http.StatusBadRequest)
			return
		}

		comment, err := m.service.PostComment.GetByID(postCommentId)
		switch {
		case err == nil:
		case errors.Is(err, model.ErrPostNotFound):
			// TODO: error Page
			http.Error(w, "Post Not Found", http.StatusNotFound)
			return
		}

		if comment.UserId != userId {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		err = m.service.PostComment.DeleteByID(comment.Id)
		switch {
		case err == nil:
		case err != nil:
			log.Printf("PostCommentDeleteHandler: m.service.PostComment.DeleteByID: %v\n", err)
			pg := &Page{Error: fmt.Errorf("internal server error, maybe try again later")}
			w.WriteHeader(http.StatusInternalServerError)
			m.tmpl.executeTemplate(w, pg, "alert.html") // TODO: Custom Epage
			return
		}

		http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
		return
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (m *mainHandler) PostCommentReactHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%-30v | %-7v | %-30v \n", r.URL, r.Method, "PostCommentReactHandler")

	switch r.Method {
	case http.MethodGet:
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	iUserId := r.Context().Value("UserId")
	if iUserId == nil {
		log.Println("PostCommentReactHandler: r.Context().Value(\"UserId\") is nil")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	userId := iUserId.(int64)

	switch r.Method {
	case http.MethodGet:
		strCommentId := r.URL.Query().Get("comment_id")
		commentId, err := strconv.ParseInt(strCommentId, 10, 64)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		strReact := r.URL.Query().Get("react")
		react, err := strconv.ParseInt(strReact, 10, 8)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		postReact := &model.PostCommentReaction{CommentId: commentId, UserId: userId, Reaction: int8(react)}
		err = m.service.PostCommentReaction.Record(postReact)
		switch {
		case err == nil:
		case errors.Is(err, model.ErrInvalidCommReaction) || errors.Is(err, model.ErrCommReactionNotFound):
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		case err != nil:
			log.Printf("PostCommentReactHandler: m.service.PostCommentReact.Record: %s", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
		return
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
}
