package routes

import (
	"net/http"
	"tiny-drop/internal/templates"

	"github.com/go-chi/chi/v5"
	// "github.com/go-chi/chi/v5/middleware"
)

func ApiRoutes(r chi.Router, renderer *templates.Renderer) {
	r.Get("/", Index(renderer))
}

func Index(renderer *templates.Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := map[string]interface{}{
			"Title": "Tiny Drop",
			"Files": []string{"example.txt"},
		}

		// Always execute "base" layout
		if err := renderer.Render(w, "base", data); err != nil {
			http.Error(w, err.Error(), 500)
		}
	}
}
