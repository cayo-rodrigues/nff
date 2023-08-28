package handlers

import (
	"html/template"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	data := map[string]bool{
		"IsAuthenticated": true,
	}
	tmpl.Execute(w, data)
}
