package main

import (
	// "html/template"
	"log"
	"net/http"
	"tiny-drop/internal/views"
	"tiny-drop/internal/routes"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)


func main() {

	render, err := views.New("web/templates")
	if err != nil {
		log.Fatalf("Error parsing templates: %v", err)
	}
	port := ":9090"
	log.Printf("Templates loaded.")
		
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	routes.ApiRoutes(r, render)
	
	
	log.Println("Server running at http://localhost", port)

	err = http.ListenAndServe(port, r)
	if err != nil {
		log.Fatal(err)
	}
}
