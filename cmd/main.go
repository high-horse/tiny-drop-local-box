package main

import (
	"log"
	"net/http"
	"tiny-drop/internal/routes"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)


func main() {
	port := ":9090"
	log.Printf("Templates loaded.")
		
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	routes.ApiRoutes(r)
	
	
	log.Printf("Server running at http://localhost%s\n", port)

	err := http.ListenAndServe(port, r)
	if err != nil {
		log.Fatal(err)
	}
}
