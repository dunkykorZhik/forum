package server

import (
	"forum/internal/config"
	"forum/internal/service"
	"net/http"
)

type mainHandler struct {
	tmpl    tmplt
	service *service.Service
}

func newMainHandler(service *service.Service, configs *config.WebCfg) (*mainHandler, error) {
	mh := &mainHandler{
		tmpl:    *newTemplate(configs.TemplatesDir),
		service: service,
	}
	return mh, nil
}

func (m *mainHandler) InitRoutes(configs *config.WebCfg) http.Handler {
	mux := http.NewServeMux()
	// HERE IS ALL ROUTES
	fsStatic := http.FileServer(http.Dir(configs.StaticFilesDir))
	mux.Handle("/static/", http.StripPrefix("/static/", fsStatic))

	// AnyRoutes
	mux.HandleFunc("/", m.IndexHandler)
	mux.HandleFunc("/signup", m.SignUpHandler)
	mux.HandleFunc("/signin", m.SignInHandler)
	mux.HandleFunc("/signout", m.SignOutHandler)

	mux.Handle("/post/get", http.HandlerFunc(m.PostViewHandler))
	mux.Handle("/post/create", m.middlewareSessionChecker(http.HandlerFunc(m.PostCreateHandler)))
	mux.Handle("/post/edit", m.middlewareSessionChecker(http.HandlerFunc(m.PostEditHandler)))
	mux.Handle("/post/react", m.middlewareSessionChecker(http.HandlerFunc(m.PostReactHandler)))
	mux.Handle("/post/delete", m.middlewareSessionChecker(http.HandlerFunc(m.PostDeleteHandler)))

	mux.Handle("/posts/own", m.middlewareSessionChecker(http.HandlerFunc(m.PostsOwnHandler)))
	mux.Handle("/posts/reacted", m.middlewareSessionChecker(http.HandlerFunc(m.PostsReactedHandler)))

	mux.Handle("/post/comment/create", m.middlewareSessionChecker(http.HandlerFunc(m.PostCommentCreateHandler)))
	mux.Handle("/post/comment/delete", m.middlewareSessionChecker(http.HandlerFunc(m.PostCommentDeleteHandler)))
	mux.Handle("/post/comment/react", m.middlewareSessionChecker(http.HandlerFunc(m.PostCommentReactHandler)))

	mux.Handle("/categories/posts", http.HandlerFunc(m.CategoriesPostsHandler))

	return mux
}
