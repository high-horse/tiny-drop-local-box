package main

import (
	"log"
	"net/http"
	"tiny-drop/internal/cleaner"
	"tiny-drop/internal/config"
	"tiny-drop/internal/db"
	"tiny-drop/internal/routes"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func init() {
	// log.LstdFlags = date + time
	// log.Lshortfile = file:line
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	db.InitDb(config.DBPath)
	port := ":9090"
	log.Printf("Templates loaded.")
		
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	routes.ApiRoutes(r)
	
	
	log.Printf("Server running at http://localhost%s\n", port)

	go cleaner.CleanupFiles()
	go cleaner.CleanupOldChunks()
	
	err := http.ListenAndServe(port, r)
	if err != nil {
		log.Fatal(err)
	}
}
