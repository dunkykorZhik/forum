package server

import (
	"errors"
	"fmt"
	"forum/internal/model"
	"log"
	"net/http"
	"time"
)

// SignUpHandler -
func (m *mainHandler) SignUpHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%-30v | %-7v | %-30v \n", r.URL, r.Method, "SignUpHandler")

	// Allowed Methods
	switch r.Method {
	case http.MethodGet:
	case http.MethodPost:
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	cookie := getSessionCookie(w, r)
	switch {
	case cookie == nil:
	case cookie != nil:
		_, err := m.service.Session.GetByUuid(cookie.Value)
		switch {
		case err == nil:
			http.Redirect(w, r, "/", http.StatusFound)
			return
		case errors.Is(err, model.ErrExpired) || errors.Is(err, model.ErrNotFound):
			addRedirectCookie(w, r.RequestURI)
			removeSessionCookie(w, r)
		case err != nil:
			log.Printf("SignUpHandler: m.service.Session.GetByUuid: %v\n", err)
			http.Error(w, "something wrong, maybe try again later", http.StatusInternalServerError)
			return
		}
	}

	// Logic
	switch r.Method {
	case http.MethodGet:
		addRedirectCookie(w, r.URL.Query().Get("redirect_to"))
		m.tmpl.executeTemplate(w, nil, "sign-up.html")
	case http.MethodPost:
		err := r.ParseForm()
		if err != nil {
			log.Printf("SignUpHandler: r.ParseForm: %v\n", err)
			pg := &Page{Error: fmt.Errorf("something wrong, maybe try again later")}
			w.WriteHeader(http.StatusInternalServerError)
			m.tmpl.executeTemplate(w, pg, "sign-in.html")
			return
		}

		newUser := &model.User{
			Username: r.FormValue("username"),
			Email:    r.FormValue("email"),
			Password: r.FormValue("password"),
		}

		_, err = m.service.User.Create(newUser)
		switch {
		case err == nil:
			http.Redirect(w, r, "/signin", http.StatusSeeOther)
			return
		case errors.Is(err, model.ErrExistUsername):
			w.WriteHeader(http.StatusBadRequest)
			pg := &Page{Error: fmt.Errorf("nickname \"%v\" is used. Try with another nickname.", newUser.Username)}
			m.tmpl.executeTemplate(w, pg, "sign-up.html")
		case errors.Is(err, model.ErrExistEmail):
			w.WriteHeader(http.StatusBadRequest)
			pg := &Page{Error: fmt.Errorf("email \"%v\" is used. Try with another email.", newUser.Email)}
			m.tmpl.executeTemplate(w, pg, "sign-up.html")
		case errors.Is(err, model.ErrInvalidUsername):
			w.WriteHeader(http.StatusBadRequest)
			pg := &Page{Error: fmt.Errorf("invalid nickname \"%v\"", newUser.Username)}
			m.tmpl.executeTemplate(w, pg, "sign-up.html")
		case errors.Is(err, model.ErrInvalidEmail):
			w.WriteHeader(http.StatusBadRequest)
			pg := &Page{Error: fmt.Errorf("invalid email \"%v\"", newUser.Email)}
			m.tmpl.executeTemplate(w, pg, "sign-up.html")
		default:
			log.Printf("SignUpHandler: %s", err)
			pg := &Page{Error: fmt.Errorf("something wrong, maybe try again later")}
			w.WriteHeader(http.StatusInternalServerError)
			m.tmpl.executeTemplate(w, pg, "sign-up.html")
			return
		}
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}
func (m *mainHandler) SignOutHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%-30v | %-7v | %-30v \n", r.URL, r.Method, "SignOutHandler")

	switch r.Method {
	case http.MethodGet:
		removeSessionCookie(w, r)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}
func (m *mainHandler) SignInHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%-30v | %-7v | %-30v \n", r.URL, r.Method, "SigInHandler")

	// Allowed Methods
	switch r.Method {
	case http.MethodGet:
	case http.MethodPost:
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	cookie := getSessionCookie(w, r)
	switch {
	case cookie == nil:
	case cookie != nil:
		_, err := m.service.Session.GetByUuid(cookie.Value)
		switch {
		case err == nil:
			http.Redirect(w, r, "/", http.StatusFound)
			return
		case errors.Is(err, model.ErrExpired) || errors.Is(err, model.ErrSessNotFound):
			addRedirectCookie(w, r.RequestURI)
			removeSessionCookie(w, r)
		case err != nil:
			log.Printf("SignInHandler: m.service.Session.GetByUuid: %v\n", err)
			http.Error(w, "something wrong, maybe try again later", http.StatusInternalServerError)
			return
		}
	}

	// Logic
	switch r.Method {
	case http.MethodGet:
		addRedirectCookie(w, r.URL.Query().Get("redirect_to"))
		m.tmpl.executeTemplate(w, nil, "sign-in.html")
		return
	case http.MethodPost:
		err := r.ParseForm()
		if err != nil {
			log.Printf("SignInHandler: r.ParseForm: %v\n", err)
			pg := &Page{Error: fmt.Errorf("something wrong, maybe try again later")}
			w.WriteHeader(http.StatusInternalServerError)
			m.tmpl.executeTemplate(w, pg, "sign-in.html")
			return
		}

		usr, err := m.service.User.GetByUsernameOrEmail(r.FormValue("login"))
		switch {
		case err == nil:
		case errors.Is(err, model.ErrNotFound):
			pg := &Page{Error: fmt.Errorf("user with login \"%v\" not found", r.FormValue("login"))}
			m.tmpl.executeTemplate(w, pg, "sign-in.html")
			return
		case errors.Is(err, model.ErrInvalidEmail):
			w.WriteHeader(http.StatusBadRequest)
			pg := &Page{Error: fmt.Errorf("invalid email %v", r.FormValue("login"))}
			m.tmpl.executeTemplate(w, pg, "sign-in.html")
			return
		case errors.Is(err, model.ErrInvalidUsername):
			w.WriteHeader(http.StatusBadRequest)
			pg := &Page{Error: fmt.Errorf("invalid nickname %v", r.FormValue("login"))}
			m.tmpl.executeTemplate(w, pg, "sign-in.html")
			return
		default:
			log.Printf("SignInHandler: User.GetByNicknameOrEmail: %s", err)
			pg := &Page{Error: fmt.Errorf("something wrong, maybe try again later")}
			w.WriteHeader(http.StatusInternalServerError)
			m.tmpl.executeTemplate(w, pg, "sign-in.html")
			return
		}

		areEqual, err := usr.CompareHashAndPassword(r.FormValue("password"))
		switch {
		case err != nil:
			log.Printf("SignInHandler: user.CompareHashAndPassword: %s", err)
			pg := &Page{Error: fmt.Errorf("something wrong, maybe try again later")}
			w.WriteHeader(http.StatusInternalServerError)
			m.tmpl.executeTemplate(w, pg, "sign-in.html")
			return
		case !areEqual:
			w.WriteHeader(http.StatusBadRequest)
			pg := &Page{Error: fmt.Errorf("invalid password for login \"%s\"", r.FormValue("login"))}
			m.tmpl.executeTemplate(w, pg, "sign-in.html")
			return
		}

		session, err := m.service.Session.Record(usr.Id)
		if err != nil {
			log.Printf("SignInHandler: Session.Record: %s", err)
			pg := &Page{Error: fmt.Errorf("something wrong, maybe try again later")}
			w.WriteHeader(http.StatusInternalServerError)
			m.tmpl.executeTemplate(w, pg, "sign-in.html")
			return
		}
		expiresAfterSeconds := time.Until(session.ExpiredAt).Seconds()
		addSessionCookie(w, session.Uuid, int(expiresAfterSeconds))

		if cookie := getRedirectCookie(w, r); cookie != nil {
			removeRedirectCookie(w, r)
			http.Redirect(w, r, cookie.Value, http.StatusFound)
			return
		}
		http.Redirect(w, r, "/", http.StatusFound)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}
