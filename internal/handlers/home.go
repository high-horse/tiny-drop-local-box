package handlers

import (
	"net/http"
	"tiny-drop/internal/views"
	"log"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	ip := query.Get("name")
	log.Println("Connected from ", ip)

	data := struct {
		Title string
		Desc  string
	}{
		Title: "sss  Home Page",
		Desc:  "some description",
	}

	views.Render(w, "layout.html", "home.html", data)
}


func ContactHandler(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Title string
		Email string
	}{
		Title: "Contact Us",
		Email: "example@example.com",
	}

	views.Render(w, "layout.html", "contact.html", data)
}

func AboutHandler(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Title string
		Desc  string
		About string
	}{
		Title: "About Us",
		Desc:  "about description",
		About: "more about about",
	}

	views.Render(w, "layout.html", "about.html", data)
}
