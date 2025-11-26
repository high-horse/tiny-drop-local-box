package views

import (
	"html/template"
	"net/http"
	"path/filepath"
	"log"
)


func Render(w http.ResponseWriter, layoutName, bodyName string, data interface{}) {
	layoutPath := filepath.Join("web", "templates", layoutName)
	bodyPath := filepath.Join("web", "templates", bodyName)

	// Parse layout and body
	layoutTmpl := template.Must(template.ParseFiles(layoutPath))
	bodyTmpl := template.Must(template.ParseFiles(bodyPath))

	// Clone layout and inject body
	tmpl := template.Must(layoutTmpl.Clone())
	tmpl = template.Must(tmpl.AddParseTree("body", bodyTmpl.Tree))

	// Execute layout template
	if err := tmpl.ExecuteTemplate(w, layoutName, data); err != nil {
		log.Printf("Template error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
