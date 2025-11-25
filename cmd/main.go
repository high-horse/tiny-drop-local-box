package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var templates *template.Template

func main() {
	templates = template.Must(template.ParseGlob("web/templates/*.html"))
	log.Printf("Parsed templates: %v", templates.DefinedTemplates())

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		render(w, "home", nil)
	})

	r.Get("/about", func(w http.ResponseWriter, r *http.Request) {
		render(w, "about", nil)
	})

	r.Get("/contact", func(w http.ResponseWriter, r *http.Request) {
		render(w, "contact", nil)
	})

	log.Println("Server running at http://localhost:9090")

	// THIS IS THE MOST IMPORTANT LINE
	err := http.ListenAndServe(":9090", r)
	if err != nil {
		log.Fatal(err)
	}
}

func render(w http.ResponseWriter, page string, data interface{}) {
	err := templates.ExecuteTemplate(w, page, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
