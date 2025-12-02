package main

import (
	"log"
	// "crypto/tls"
	"net/http"
	"tiny-drop/internal/cleaner"
	"tiny-drop/internal/config"
	"tiny-drop/internal/db"
	"tiny-drop/internal/routes"
	"tiny-drop/internal/views"

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
	port := config.Port
	views.InitTemplates()
	log.Printf("Templates loaded.")
		
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	routes.ApiRoutes(r)
	
	
	go cleaner.StartCleanupTicker()
	go cleaner.CleanupOldChunks()
	
	// for server
	// log.Printf("Server running at http://localhost%s\n", port)
	// err := http.ListenAndServe(port, r)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// for local
	srv := &http.Server{
		Addr: port,
		Handler: r,
	}

	log.Printf("Server running at https://localhost%s\n", port)
	err := srv.ListenAndServeTLS("server.crt", "server.key")
	if err != nil {
		log.Fatal(err)
	}
}
