package handlers

import (
	"net/http"
	"tiny-drop/internal/services"
	"tiny-drop/internal/utils"
	"tiny-drop/internal/views"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	ip := query.Get("ip")
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
		"Title" : "Tiny Drop – Fast & Easy File Sharing on Local & Public Networks",
		"Desc" : "Tiny Drop makes sharing files simple, secure, and lightning-fast. Instantly send and receive files across local networks or over the internet without hassle. Perfect for personal, team, or public use.",
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
