package routes

import (
	"log"
	"net/http"
	"tiny-drop/internal/views"

	"github.com/go-chi/chi/v5"
	// "github.com/go-chi/chi/v5/middleware"
)

func ApiRoutes(r chi.Router) {
	r.Get("/", homeHandler)
	r.Get("/about", aboutHandler)
	r.Get("/contact", contactHandler)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	ip := query.Get("name")
	log.Println("Connected from ", ip)
	
	data := struct {
		Title string
		Desc string
	}{
		Title: "sss  Home Page",
		Desc: "some description",
	}

	views.Render(w, "new_layout.html", "home_body.html", data)
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Title string
		Email string
	}{
		Title: "Contact Us",
		Email: "example@example.com",
	}

	views.Render(w, "new_layout.html", "contact_body.html", data)
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Title string
		Desc string
		About string
	}{
		Title: "About Us",
		Desc: "about description",
		About: "more about about",
	}

	views.Render(w, "new_layout.html", "about_body.html", data)
}


