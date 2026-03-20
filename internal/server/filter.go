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

func (m *mainHandler) PostsReactedHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%-30v | %-7v | %-30v \n", r.URL, r.Method, "PostsReactedHandler")

	// Allowed Methods
	switch r.Method {
	case http.MethodGet:
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	iUserId := r.Context().Value("UserId")
	if iUserId == nil {
		log.Println("PostsReactedHandler: r.Context().Value(\"UserId\") is nil")
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
		log.Printf("PostsReactedHandler: m.service.User.GetByID: %v\n", err)
		pg := &Page{Error: fmt.Errorf("internal server error, maybe try again later")}
		w.WriteHeader(http.StatusInternalServerError)
		m.tmpl.executeTemplate(w, pg, "alert.html")
		return
	}

	switch r.Method {
	case http.MethodGet:
		strReact := r.URL.Query().Get("react")
		react, err := strconv.ParseInt(strReact, 10, 8)
		if err != nil || react < -1 || 1 < react {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		postIDs, err := m.service.PostReaction.GetAllUserReactedPostIDs(userId, int8(react), 0, 0)
		if err != nil {
			log.Printf("PostsReactedHandler: PostVote.GetAllUserReactedPostIDs: %v\n", err)
			pg := &Page{Error: fmt.Errorf("internal server error, maybe try again later")}
			w.WriteHeader(http.StatusInternalServerError)
			m.tmpl.executeTemplate(w, pg, "alert.html") // TODO: Custom Epage
			return
		}

		posts, err := m.service.Post.GetByIDs(postIDs)
		if err != nil {
			log.Printf("PostsReactedHandler: Post.GetByIDs: %v\n", err)
		}

		err = m.service.FillPosts(posts, user.Id)
		if err != nil {
			log.Printf("PostsReactedHandler: FillPosts: %v\n", err)
		}

		pg := &Page{User: user, Posts: posts, Info: fmt.Errorf("Reacted Posts")}
		m.tmpl.executeTemplate(w, pg, "home.html")
		return
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
}

func (m *mainHandler) PostsOwnHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%-30v | %-7v | %-30v \n", r.URL, r.Method, "PostsOwnHandler")
	// Allowed Methods
	switch r.Method {
	case http.MethodGet:
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	iUserId := r.Context().Value("UserId")
	if iUserId == nil {
		log.Println("PostsOwnHandler: r.Context().Value(\"UserId\") is nil")
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
		log.Printf("PostsOwnHandler: m.service.User.GetByID: %v\n", err)
		pg := &Page{Error: fmt.Errorf("internal server error, maybe try again later")}
		w.WriteHeader(http.StatusInternalServerError)
		m.tmpl.executeTemplate(w, pg, "alert.html") // TODO: Custom Epage
		return
	}

	switch r.Method {
	case http.MethodGet:
		posts, err := m.service.Post.GetByUserID(user.Id, 0, 0)
		if err != nil {
			log.Printf("PostsOwnHandler: GetByUserID: %v\n", err)
		}

		err = m.service.FillPosts(posts, user.Id)
		if err != nil {
			log.Printf("PostsOwnHandler: FillPosts: %v\n", err)
		}

		pg := &Page{User: user, Posts: posts, Info: fmt.Errorf("Here is your posts")}
		m.tmpl.executeTemplate(w, pg, "home.html")
		return
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
}

func (m *mainHandler) CategoriesPostsHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%-30v | %-7v | %-30v \n", r.URL, r.Method, "PostsCategoriesHandler")
	switch r.Method {
	case http.MethodGet:
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	switch r.Method {
	case http.MethodGet:
		value := r.URL.Query().Get("categories")
		categoryNames := strings.Fields(value)
		if len(categoryNames) > 5 {
			http.Error(w, "Max Category names is 5", http.StatusBadRequest)
			return
		}

		categories, err := m.service.Category.GetByNames(categoryNames)
		switch {
		case err == nil:
		case err != nil:
			log.Printf("CategoriesPostsHandler: Category.GetByNames: %v\n", err)
			http.Error(w, "something wrong, maybe try again later", http.StatusInternalServerError)
			return
		}

		var infoMsg error
		if len(categories) != len(categoryNames) {
			infoMsg = fmt.Errorf("Looking for only contaned categories")
		}

		catIDs := make([]int64, len(categories))
		for i, v := range categories {
			catIDs[i] = v.Id
		}

		postIDs, err := m.service.Category.GetPostIDsContainedCatIDs(catIDs, 0, -1)
		switch {
		case err == nil:
		case err != nil:
			log.Printf("CategoriesPostsHandler: Category.GetPostIDsContainedCatIDs: %v\n", err)
			http.Error(w, "something wrong, maybe try again later", http.StatusInternalServerError)
			return
		}

		// TODO: Rename Ids -> IDs
		posts, err := m.service.Post.GetByIDs(postIDs)
		switch {
		case err == nil:
		case err != nil:
			log.Printf("CategoriesPostsHandler: Post.GetByIDs: %v\n", err)
			http.Error(w, "something wrong, maybe try again later", http.StatusInternalServerError)
			return
		}

		pg := &Page{Posts: posts, Categories: categories, Info: infoMsg}
		cookie := getSessionCookie(w, r)
		if cookie == nil {
			err = m.service.FillPosts(posts, 0)
			if err != nil {
				log.Printf("CategoriesPostsHandler: FillPosts: %v\n", err)
			}
			m.tmpl.executeTemplate(w, pg, "categories-posts.html")
			return
		}

		session, err := m.service.Session.GetByUuid(cookie.Value)
		switch {
		case err == nil:
		case errors.Is(err, model.ErrExpired) || errors.Is(err, model.ErrSessNotFound):
			removeSessionCookie(w, r)
			err = m.service.FillPosts(posts, 0)
			if err != nil {
				log.Printf("CategoriesPostsHandler: FillPosts: %v\n", err)
			}
			m.tmpl.executeTemplate(w, pg, "categories-posts.html")
			return
		case err != nil:
			log.Printf("CategoriesPostsHandler: m.service.Session.GetByUuid: %v\n", err)
			http.Error(w, "something wrong, maybe try again later", http.StatusInternalServerError)
			return
		}

		user, err := m.service.User.GetByID(session.UserId)
		switch {
		case err == nil:
		case err != nil:
			log.Printf("CategoriesPostsHandler: m.service.Session.GetByUuid: %v\n", err)
			http.Error(w, "something wrong, maybe try again later", http.StatusInternalServerError)
			return
		}

		err = m.service.FillPosts(posts, user.Id)
		if err != nil {
			log.Printf("CategoriesPostsHandler: FillPosts: %v\n", err)
		}
		pg.User = user
		m.tmpl.executeTemplate(w, pg, "categories-posts.html")
		return
	}
}
