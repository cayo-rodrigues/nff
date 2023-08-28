package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/cayo-rodrigues/nff/web/handlers"
)

func main() {
	http.HandleFunc("/static/styles/", handlers.ServeStyles)

	http.HandleFunc("/", handlers.Index)
	http.HandleFunc("/entities", handlers.Entities)

	fmt.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
