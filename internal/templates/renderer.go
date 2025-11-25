package templates

import (
	"html/template"
	"net/http"
	"path/filepath"
)

type Renderer struct {
	templates *template.Template
	basePath  string
	debug     bool
}

func NewRenderer(basePath string, debug bool) *Renderer {
	r := &Renderer{
		basePath: basePath,
		debug:    debug,
	}
	r.ParseTemplates()
	return r
}

func (r *Renderer) ParseTemplates() {
	pattern := filepath.Join(r.basePath, "*.tmpl")
	r.templates = template.Must(template.ParseGlob(pattern))
}

func (r *Renderer) Render(w http.ResponseWriter, name string, data interface{}) error {
	// auto-reload in debug mode
	if r.debug {
		r.ParseTemplates()
	}

	return r.templates.ExecuteTemplate(w, name, data)
}
