package routes

import (
	"html/template"
	"net/http"
	"tiny-drop/internal/views"

	"github.com/go-chi/chi/v5"
	"path/filepath"
	"log"
	// "github.com/go-chi/chi/v5/middleware"
)

func ApiRoutes(r chi.Router, render *views.Renderer) {
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		render.Render(w, "home", nil)
	})

	r.Get("/about", func(w http.ResponseWriter, r *http.Request) {
		render.Render(w, "about", nil)
	})

	// r.Get("/contact", func(w http.ResponseWriter, r *http.Request) {
	// 	data := map[string]any{
	// 		"Title": "Contact Us",
	// 		"Email": "example@example.com",
	// 	}
	// 	render.Render(w, "contact", data)
	// })
	// 
	// r.Get("/contact", func(w http.ResponseWriter, r *http.Request) {
	// 	tmpl := template.Must(template.ParseFiles("web/templates/new_content.html"))
	// 	data := map[string]any{
	// 		"Title": "Contact Us",
	// 		"Email": "example@example.com",
	// 	}
	// 	if err := tmpl.Execute(w, data); err != nil {
	// 		http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	}
	// })
	// 
	// r.Get("/contact", func(w http.ResponseWriter, r *http.Request) {
	//     layout := template.Must(template.ParseFiles("web/templates/new_layout.html"))
	//     body   := template.Must(template.ParseFiles("web/templates/contact_body.html"))
	
	//     // Inject body as named template "body"
	//     tmpl := template.Must(layout.Clone())
	//     tmpl = template.Must(tmpl.AddParseTree("body", body.Tree))
	
	//     data := map[string]any{
	//         "Title": "Contact Us",
	//         "Email": "example@example.com",
	//     }
	
	//     tmpl.ExecuteTemplate(w, "layout.html", data)
	// })
	r.Get("/contact", contactHandler)

}


func contactHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Parse layout and body templates
	layoutPath := filepath.Join("web", "templates", "new_layout.html")
	bodyPath   := filepath.Join("web", "templates", "contact_body.html")

	layoutTmpl := template.Must(template.ParseFiles(layoutPath))
	bodyTmpl   := template.Must(template.ParseFiles(bodyPath))

	// 2. Clone layout and inject body as named template "body"
	tmpl := template.Must(layoutTmpl.Clone())
	tmpl = template.Must(tmpl.AddParseTree("body", bodyTmpl.Tree))

	// 3. Prepare data
	data := struct {
		Title string
		Email string
	}{
		Title: "Contact Us",
		Email: "example@example.com",
	}

	// 4. Render layout (which calls {{ template "body" . }})
	err := tmpl.ExecuteTemplate(w, "new_layout.html", data)
	if err != nil {
		log.Printf("Template error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}