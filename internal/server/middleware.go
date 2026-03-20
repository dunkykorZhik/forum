package server

import (
	"context"
	"errors"
	"forum/internal/model"
	"log"
	"net/http"
)

// MiddlewareSessionChecker - NOT FINISHED
func (m *mainHandler) middlewareSessionChecker(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%-30v | %-7v | %-30v \n", r.URL, r.Method, "middlewareSessionChecker")
		cookie := getSessionCookie(w, r)
		if cookie == nil {
			if r.Method == http.MethodGet {
				addRedirectCookie(w, r.RequestURI)
			}
			http.Redirect(w, r, "/signin", http.StatusSeeOther)
			return
		}

		session, err := m.service.Session.GetByUuid(cookie.Value)
		switch {
		case err == nil:
		case errors.Is(err, model.ErrExpired) || errors.Is(err, model.ErrNotFound):
			if r.Method == http.MethodGet {
				addRedirectCookie(w, r.RequestURI)
			}
			removeSessionCookie(w, r)
			http.Redirect(w, r, "/signin", http.StatusSeeOther)
			return
		case err != nil:
			log.Printf("MiddlewareSessionChecker: m.service.Session.GetByUuid: %v\n", err)
			http.Error(w, "something wrong, maybe try again later", http.StatusInternalServerError)
			return
		}
		ctx := r.Context()
		ctx = context.WithValue(ctx, "UserId", session.UserId)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

/*
func (m *mainHandler) middlewareMethodChecker(next http.Handler, allowedMthods map[string]bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%-30v | %-7v | %-30v \n", r.URL, r.Method, "middlewareMethodChecker")

		if _, ok := allowedMthods[r.Method]; !ok {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
		next.ServeHTTP(w, r)
	})
}

*/
