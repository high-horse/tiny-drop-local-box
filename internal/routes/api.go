package routes

import (
	"tiny-drop/internal/handlers"

	"github.com/go-chi/chi/v5"
	// "github.com/go-chi/chi/v5/middleware"
)

func ApiRoutes(r chi.Router) {
	r.Get("/", handlers.HomeHandler)
	r.Get("/about", handlers.AboutHandler)
	r.Get("/contact", handlers.ContactHandler)
	
	r.Post("/upload", handlers.UploadHandler)
}
