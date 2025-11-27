package handlers

import (
	"net/http"
	"tiny-drop/internal/services"
	"tiny-drop/internal/utils"
	"tiny-drop/internal/views"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	ip := query.Get("name")
	if ip == "" {
		ip = utils.GetUserIp(r)
	}

	files, err := services.FetchUploadeds(ip)
	if err != nil {
		views.Render(w, "layout.html", "home.html", map[string]any {
			"success" : false,
			"message" :  "failed to fetch files",
		})
		return
	}

	data := map[string]any {
		"Title" : "Home Page",
		"Desc" : "Home Page Description",
		"Data" : files,
		"Success": true,
		"Ip" : ip,
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
