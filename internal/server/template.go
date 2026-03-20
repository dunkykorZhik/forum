package server

import (
	"fmt"
	"forum/internal/model"
	"html/template"
	"log"
	"net/http"
)

type tmplt struct {
	templatesDir string
}

type Page struct {
	User       *model.User
	Users      []*model.User
	Post       *model.Post
	Posts      []*model.Post
	Categories []*model.Category

	// Comments           []models.Comment
	Error   error // Error - Notification Error
	Warn    error // Warn - Notification Warning
	Info    error // Info - Notification Info
	Success error // Success - Notification Success
}

func newTemplate(templatesDir string) *tmplt {
	return &tmplt{templatesDir: templatesDir}
}

func (t tmplt) executeTemplate(w http.ResponseWriter, pg interface{}, names ...string) {
	tmpl, err := t.getTemplate(names...)
	if err != nil {
		log.Printf("m.newView: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "main", pg)
	if err != nil {
		log.Printf("tmpl.ExecuteTemplate: %v", err)
		return
	}
}

func (t *tmplt) getTemplate(names ...string) (*template.Template, error) {
	paths := []string{t.templatesDir + "/main.html", t.templatesDir + "/navbar.html", t.templatesDir + "/alert.html"}
	for _, name := range names {
		paths = append(paths, t.templatesDir+"/"+name)
	}

	res, err := template.ParseFiles(paths...)
	if err != nil {
		return nil, fmt.Errorf("tmplt.ParseFiles: %w", err)
	}
	return res, nil
}
