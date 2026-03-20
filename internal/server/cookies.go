package server

import (
	"errors"
	"log"
	"net/http"
)

const (
	CookieSessionName  = "session"
	CookieRedirectName = "redirect_to"
)

func addRedirectCookie(w http.ResponseWriter, redirectTo string) {
	if redirectTo == "" {
		return
	}
	http.SetCookie(w,
		&http.Cookie{
			Name:   CookieRedirectName,
			Value:  redirectTo,
			Path:   "/",
			MaxAge: 3600,
		},
	)
}

func getRedirectCookie(w http.ResponseWriter, r *http.Request) *http.Cookie {
	return getCookieByName(w, r, CookieRedirectName)
}

func removeRedirectCookie(w http.ResponseWriter, r *http.Request) {
	removeCookieByName(w, r, CookieRedirectName)
}

func addSessionCookie(w http.ResponseWriter, uuid string, durationSec int) {
	http.SetCookie(w,
		&http.Cookie{
			Name:   CookieSessionName,
			Value:  uuid,
			Path:   "/",
			MaxAge: durationSec,
		},
	)
}

func getSessionCookie(w http.ResponseWriter, r *http.Request) *http.Cookie {
	return getCookieByName(w, r, CookieSessionName)
}

// RemoveSessionCookie - removes cookie by setting maxAge -1
func removeSessionCookie(w http.ResponseWriter, r *http.Request) {
	removeCookieByName(w, r, CookieSessionName)
}

func getCookieByName(w http.ResponseWriter, r *http.Request, name string) *http.Cookie {
	cookie, err := r.Cookie(name)
	switch {
	case errors.Is(err, http.ErrNoCookie):
	case err != nil:
		log.Printf("GetRedirectCookie: r.Cookie: %v", err)
	case cookie != nil:
		return cookie
	}
	return nil
}

// removeCookieByName - remove cookie by setting maxAge -1
func removeCookieByName(w http.ResponseWriter, r *http.Request, name string) {
	cookie, err := r.Cookie(name)
	switch {
	case errors.Is(err, http.ErrNoCookie):
	case err != nil:
		log.Printf("removeCookieByName: r.Cookie: %v", err)
	case cookie != nil:
		cookie.MaxAge = -1
		http.SetCookie(w, cookie)
	}
}
