package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

func ServeStyles(w http.ResponseWriter, r *http.Request) {
	_, stylesheet := filepath.Split(r.URL.Path)

	filepath := fmt.Sprintf("static/styles/%s", stylesheet)
	_, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		fmt.Printf("File '%s' does not exist.\n", filepath)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	http.ServeFile(w, r, filepath)
}
