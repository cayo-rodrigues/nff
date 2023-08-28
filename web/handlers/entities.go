package handlers

import (
	"html/template"
	"net/http"
)

func Entities(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		getEntities(w, r)
	}
}

func getEntities(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(
		"templates/layout.html", "templates/entities.html",
	))
	tmpl.ExecuteTemplate(w, "layout", nil)
}
