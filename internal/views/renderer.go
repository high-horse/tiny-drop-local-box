package views

import (
	"html/template"
	"net/http"
	"path/filepath"
	"log"
	"sync"
)

var (
    templates = make(map[string]*template.Template)
    once      sync.Once
)

func InitTemplates() {
    once.Do(func() {
        layoutPath := filepath.Join("web", "templates", "layout.html")
        
        // Parse all page templates
        pages := []string{"home.html"} // add more pages as needed
        
        for _, page := range pages {
            pagePath := filepath.Join("web", "templates", page)
            
            // Parse both files together - this will work now
            tmpl, err := template.ParseFiles(layoutPath, pagePath)
            if err != nil {
                log.Fatalf("Error parsing template %s: %v", page, err)
            }
            
            templates[page] = tmpl
        }
        
        log.Printf("Loaded %d templates", len(templates))
    })
}

func InitTemplates_() {
    once.Do(func() {
        layoutPath := filepath.Join("web", "templates", "layout.html")
        
        // Parse all page templates
        pages := []string{"home.html"} // add more pages as needed
        
        for _, page := range pages {
            pagePath := filepath.Join("web", "templates", page)
            
            // Parse layout + page together
            tmpl, err := template.ParseFiles(layoutPath, pagePath)
            if err != nil {
                log.Fatalf("Error parsing template %s: %v", page, err)
            }
            
            templates[page] = tmpl
        }
        
        log.Printf("Loaded %d templates", len(templates))
    })
}


func Render(w http.ResponseWriter, layoutName, bodyName string, data interface{}) {
    tmpl, ok := templates[bodyName]
    if !ok {
        log.Printf("Template %s not found", bodyName)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }
    log.Println(layoutName, bodyName)
    if err := tmpl.ExecuteTemplate(w, layoutName, data); err != nil {
        log.Printf("Template execution error: %v", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
    }
}

func RenderOld(w http.ResponseWriter, layoutName, bodyName string, data interface{}) {
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
