package server

import (
	"errors"
	"fmt"
	"forum/internal/model"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func (m *mainHandler) PostCreateHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%-30v | %-7v | %-30v \n", r.URL, r.Method, "PostCreateHandler")

	// Allowed Methods
	switch r.Method {
	case http.MethodGet:
	case http.MethodPost:
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	iUserId := r.Context().Value("UserId")
	if iUserId == nil {
		log.Println("PostCreateHandler: r.Context().Value(\"UserId\") is nil")
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
		addRedirectCookie(w, r.RequestURI)
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
	case http.MethodGet:
		pg := &Page{User: user}
		m.tmpl.executeTemplate(w, pg, "post-create.html")
		return
	case http.MethodPost:
		r.ParseForm()

		post := &model.Post{
			Title:   r.FormValue("title"),
			Content: r.FormValue("content"),
			UserId:  userId,
		}
		_, err := m.service.Post.Create(post)
		switch {
		case err == nil:
		case errors.Is(err, model.ErrInvalidTitleLength):
			w.WriteHeader(http.StatusBadRequest)
			pg := &Page{Error: fmt.Errorf("invalid length of title")}
			m.tmpl.executeTemplate(w, pg, "post-create.html")
			return
		case errors.Is(err, model.ErrInvalidContentLength):
			w.WriteHeader(http.StatusBadRequest)
			pg := &Page{Error: fmt.Errorf("invalid length of content")}
			m.tmpl.executeTemplate(w, pg, "post-create.html")
			return
		default:
			log.Printf("PostCreateHandler: m.service.Post.Create: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			pg := &Page{Error: fmt.Errorf("something wrong, maybe try again later: %s", err)}
			m.tmpl.executeTemplate(w, pg, "post-create.html")
			return
		}

		catNames := strings.Fields(r.Form.Get("categories"))
		err = m.service.Category.AddToPostByNames(catNames, post.Id)
		switch {
		case err == nil:
		case errors.Is(err, model.ErrCategoryLimitForPost):
			err = m.service.Post.DeleteByID(post.Id)
			if err != nil {
				log.Println("PostCreateHandler: m.service.Post.DeleteByID: %w", err)
			}

			w.WriteHeader(http.StatusBadRequest)
			pg := &Page{Error: fmt.Errorf("post not created, invalid categies count, category limit = %v", 5)}
			m.tmpl.executeTemplate(w, pg, "post-create.html")
			return
		default:
			err = m.service.Post.DeleteByID(post.Id)
			if err != nil {
				log.Println("PostCreateHandler: m.service.Post.DeleteByID: %w", err)
			}

			log.Printf("PostCreateHandler:  m.service.Category.AddToPostByNames: %s", err)
			pg := &Page{Error: fmt.Errorf("something wrong, maybe try again later: %s", err)}
			w.WriteHeader(http.StatusInternalServerError)
			m.tmpl.executeTemplate(w, pg, "post-create.html")
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/post/get?id=%v", post.Id), http.StatusSeeOther)
		return
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (m *mainHandler) PostDeleteHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%-30v | %-7v | %-30v \n", r.URL, r.Method, "PostDeleteHandler")

	switch r.Method {
	case http.MethodGet:
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	iUserId := r.Context().Value("UserId")
	if iUserId == nil {
		log.Println("PostDeleteHandler: r.Context().Value(\"UserId\") is nil")
		pg := &Page{Error: fmt.Errorf("internal server error, maybe try again later")}
		w.WriteHeader(http.StatusInternalServerError)
		m.tmpl.executeTemplate(w, pg, "alert.html") // TODO: Custom Epage
		return
	}

	userId := iUserId.(int64)
	user, err := m.service.User.GetByID(userId)
	switch {
	case err == nil:
	case errors.Is(err, model.ErrNotFound):
		removeSessionCookie(w, r)
		addRedirectCookie(w, r.RequestURI)
		http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
		return
	case err != nil:
		log.Printf("PostDeleteHandler: m.service.User.GetByID: %v\n", err)
		pg := &Page{Error: fmt.Errorf("internal server error, maybe try again later")}
		w.WriteHeader(http.StatusInternalServerError)
		m.tmpl.executeTemplate(w, pg, "alert.html") // TODO: Custom Epage
		return
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
		case errors.Is(err, model.ErrNotFound):
			// TODO: error Page
			http.Error(w, "Post Not Found", http.StatusNotFound)
			return
		}

		if post.UserId != user.Id {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		err = m.service.Post.DeleteByID(post.Id)
		switch {
		case err == nil:
		case err != nil:
			log.Printf("PostDeleteHandler: m.service.Post.DeleteByID: %v\n", err)
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

func (m *mainHandler) PostEditHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%-30v | %-7v | %-30v \n", r.URL, r.Method, "PostEditHandler")

	// Allowed Methods
	switch r.Method {
	case http.MethodGet:
	case http.MethodPost:
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	iUserId := r.Context().Value("UserId")
	if iUserId == nil {
		log.Println("PostEditHandler: r.Context().Value(\"UserId\") is nil")
		pg := &Page{Error: fmt.Errorf("internal server error, maybe try again later")}
		w.WriteHeader(http.StatusInternalServerError)
		m.tmpl.executeTemplate(w, pg, "post-edit.html")
		return
	}

	userId := iUserId.(int64)
	user, err := m.service.User.GetByID(userId)
	switch {
	case err == nil:
	case errors.Is(err, model.ErrNotFound):
		removeSessionCookie(w, r)
		addRedirectCookie(w, r.RequestURI)
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
		case errors.Is(err, model.ErrNotFound):
			// TODO: error Page
			http.Error(w, "Post Not Found", http.StatusNotFound)
			return
		}

		if post.UserId != user.Id {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		categories, err := m.service.Category.GetByPostID(post.Id)
		switch {
		case err == nil:
			post.Categories = categories
		default:
			log.Printf("PostEditHandler: m.service.PostCategory.GetByPostID: %v\n", err)
			http.Error(w, "something wrong, maybe try again later", http.StatusInternalServerError)
			return
		}

		pg := &Page{User: user, Post: post}
		m.tmpl.executeTemplate(w, pg, "post-edit.html")
		return
	case http.MethodPost:
		r.ParseForm()

		var strPostId string = r.FormValue("id")
		postId, err := strconv.ParseInt(strPostId, 10, 64)
		if err != nil || postId < 1 {
			http.Error(w, "Invalid query id", http.StatusBadRequest)
			return
		}

		post, err := m.service.Post.GetByID(postId)
		switch {
		case err == nil:
		case errors.Is(err, model.ErrNotFound):
			// TODO: error Page
			http.Error(w, "Post Not Found", http.StatusNotFound)
			return
		}

		if post.UserId != user.Id {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		post = &model.Post{
			Id:      postId,
			Title:   r.FormValue("title"),
			Content: r.FormValue("content"),
			UserId:  user.Id,
		}
		err = m.service.Post.Update(post)
		switch {
		case err == nil:
		case errors.Is(err, model.ErrInvalidTitleLength) || errors.Is(err, model.ErrInvalidContentLength):
			categories, errn := m.service.Category.GetByPostID(post.Id)
			switch {
			case errn == nil:
				post.Categories = categories
			default:
				log.Printf("PostEditHandler: m.service.PostCategory.GetByPostID: %v\n", err)
				http.Error(w, "something wrong, maybe try again later", http.StatusInternalServerError)
				return
			}

			var errMsg error
			switch {
			case errors.Is(err, model.ErrInvalidTitleLength):
				errMsg = fmt.Errorf("invalid title length")
			case errors.Is(err, model.ErrInvalidContentLength):
				errMsg = fmt.Errorf("invalid content length")
			default:
				log.Printf("PostEditHandler: havent got message for error: %s\n", err)
				errMsg = fmt.Errorf("invalid post")
			}

			w.WriteHeader(http.StatusBadRequest)
			pg := &Page{User: user, Post: post, Error: errMsg}
			m.tmpl.executeTemplate(w, pg, "post-edit.html")
			return
		case err != nil:
			log.Printf("PostEditHandler: m.service.Post.Update: %v\n", err)
			http.Error(w, "something wrong, maybe try again later", http.StatusInternalServerError)
			return
		}

		err = m.service.Category.DeleteByPostID(post.Id)
		switch {
		case err == nil:
		case err != nil:
			log.Printf("PostEditHandler: m.service.PostCategory.DeleteByPostID: %v\n", err)
			http.Error(w, "something wrong, maybe try again later", http.StatusInternalServerError)
			return
		}

		catNames := strings.Fields(r.Form.Get("categories"))
		err = m.service.Category.AddToPostByNames(catNames, post.Id)
		switch {
		case err == nil:
		case errors.Is(err, model.ErrCategoryLimitForPost):
			pg := &Page{Warn: fmt.Errorf("post categories not updated, invalid categies count, category limit = %v", 5), Post: post}
			w.WriteHeader(http.StatusBadRequest)
			m.tmpl.executeTemplate(w, pg, "post-edit.html")
			return
		default:
			log.Printf("PostEditHandler:  m.service.Category.AddToPostByNames: %s", err)
			pg := &Page{Error: fmt.Errorf("something wrong, maybe try again later: %s", err), Post: post}
			w.WriteHeader(http.StatusInternalServerError)
			m.tmpl.executeTemplate(w, pg, "post-edit.html")
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/post/get?id=%v", post.Id), http.StatusSeeOther)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (m *mainHandler) PostReactHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%-30v | %-7v | %-30v \n", r.URL, r.Method, "PostReactHandler")

	switch r.Method {
	case http.MethodGet:
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	iUserId := r.Context().Value("UserId")
	if iUserId == nil {
		log.Println("PostReactHandler: r.Context().Value(\"UserId\") is nil")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	userId := iUserId.(int64)

	switch r.Method {
	case http.MethodGet:
		strPostId := r.URL.Query().Get("post_id")
		postId, err := strconv.ParseInt(strPostId, 10, 64)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		strVote := r.URL.Query().Get("react")
		vote, err := strconv.ParseInt(strVote, 10, 8)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		postReaction := &model.PostReaction{PostId: postId, UserId: userId, Reaction: int8(vote)}
		err = m.service.PostReaction.Record(postReaction)
		switch {
		case err == nil:
		case errors.Is(err, model.ErrInvalidPostReaction) || errors.Is(err, model.ErrPostReactNotFound):
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		case err != nil:
			log.Printf("PostVoteHandler: m.service.PostVote.Record: %s", err)
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
