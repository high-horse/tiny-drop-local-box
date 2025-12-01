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

	r.Get("/download", handlers.HandleDownloadStream)
	r.Get("/file-info", handlers.HandleFileInfo)
	r.Get("/delete", handlers.DeleteHandler)

	r.Get("/events", handlers.SSEHandler)

}
