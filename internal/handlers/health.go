package handlers

import (
	"net/http"
	"tiny-drop/internal/views"
)

func HandleHealth(w http.ResponseWriter, r *http.Request) {
	data := map[string]any{
		"Title": "Contact Us",
	}
	views.GlobalRenderer.Render(w, "health", data)
}