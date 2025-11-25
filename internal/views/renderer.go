package views

import (
	"html/template"
	"net/http"
	"path/filepath"
)

type Renderer struct {
	templates *template.Template
}

var GlobalRenderer Renderer

func New(dir string) (*Renderer, error) {
	tmpl, err := template.ParseGlob(filepath.Join(dir, "*.html"))
	if err != nil {
		return nil, err
	}

	GlobalRenderer = Renderer{
		templates: tmpl,
	}
	return &GlobalRenderer, nil
}

func (r *Renderer) Render(w http.ResponseWriter, name string, data any) {
	err := r.templates.ExecuteTemplate(w, name, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
