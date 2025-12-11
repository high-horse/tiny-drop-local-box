package db

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	// _ "github.com/mattn/go-sqlite3"
	_ "modernc.org/sqlite"
)


var(
	db *sql.DB
	once sync.Once
)


func InitDb(dataSourceName string){
	once.Do(func() {
		var err error
		// db, err = sql.Open("sqlite3", dataSourceName)
		db, err = sql.Open("sqlite", dataSourceName)

		if err != nil {
			log.Fatalf("failed to initialize database: %v ", err)
		}

		err = db.Ping()
		if err != nil {
			log.Fatalf("failed to connect to database : %v", err)
		}

		fmt.Println("Database connection extablished successfully.")
		CreateTable()
	})
}


func GetDB() *sql.DB {
	if db == nil {
		log.Fatal("Database not initialized, call INIT first.")
	}
	return db
}

func CreateTable() {
	db := GetDB()

	// PRAGMAs (all valid for modernc.org/sqlite)
	if _, err := db.Exec("PRAGMA journal_mode = WAL;"); err != nil {
		log.Fatalf("Failed to create WAL mode :%v", err)
	}

	 // Optimize SQLite performance
    if _, err := db.Exec("PRAGMA synchronous = NORMAL;"); err != nil {
        log.Printf("Warning: Failed to set synchronous mode: %v", err)
    }
    
    if _, err := db.Exec("PRAGMA cache_size = -64000;"); err != nil { // 64MB cache
        log.Printf("Warning: Failed to set cache size: %v", err)
    }
    
    if _, err := db.Exec("PRAGMA temp_store = MEMORY;"); err != nil {
        log.Printf("Warning: Failed to set temp_store: %v", err)
    }

	createTableSql := `
	CREATE TABLE IF NOT EXISTS uploads (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		ip TEXT,
		uploader_id TEXT,
		file_uuid TEXT,
		file_name TEXT NOT NULL,
		file_path TEXT NOT NULL,
		file_size INTEGER NOT NULL,
		uploaded_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		last_download_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		metadata TEXT
	);
	`
	_, err := db.Exec(createTableSql)
	if err != nil {
		log.Fatalf("Failed to create Table: %v", err)
	}
}