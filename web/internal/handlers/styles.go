package handlers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

func ServeStyles(w http.ResponseWriter, r *http.Request) {
	stylesheet := chi.URLParam(r, "stylesheet")

	filepath := fmt.Sprintf("internal/static/styles/%s", stylesheet)
	_, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		fmt.Printf("File '%s' does not exist.\n", filepath)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	http.ServeFile(w, r, filepath)
}
