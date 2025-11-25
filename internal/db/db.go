package db

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)


var(
	db *sql.DB
	once sync.Once
)


func InitDb(dataSourceName string){
	once.Do(func() {
		var err error
		db, err = sql.Open("sqlite3", dataSourceName)
		if err != nil {
			log.Fatalf("failed to initialize database: %v ", err)
		}

		err = db.Ping()
		if err != nil {
			log.Fatalf("failed to connect to database : %v", err)
		}

		fmt.Println("Database connection extablished successfully.")
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
	createTableSql := `
	CREATE TABLE IF NOT EXISTS uploads (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		uuid TEXT,
		filename TEXT NOT NULL,
		filepath TEXT NOT NULL,
		filesize INTEGER NOT NULL,
		uploaded_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		metadata TEXT
	);
	`
	_, err := db.Exec(createTableSql)
	if err != nil {
		log.Fatalf("Failed to create Table: %v", err)
	}
}