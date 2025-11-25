package main

import (
	"log"
	"net/http"
	"tiny-drop/internal/routes"
	"tiny-drop/internal/templates"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	
	fileServer := http.FileServer(http.Dir("./web/static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fileServer))
	
	renderer := templates.NewRenderer("./web/templates", true)
	
	
	routes.ApiRoutes(r, renderer)
	
	port := ":3000"
	log.Println("Listening on ", port)
	http.ListenAndServe(port, r)
}