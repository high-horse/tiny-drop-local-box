package cleaner

import (
	"log"
	"os"
	"path/filepath"
	"time"
	"tiny-drop/internal/config"
	"tiny-drop/internal/db"
)

const (
	tempDir  = config.TempDir
	// finalDir = config.FinalDir
)

func CleanupOldChunks() {
	ticker := time.NewTicker(config.CleanupTime)
	defer ticker.Stop()

	for {
		<-ticker.C // wait for next tick

		err := filepath.Walk(tempDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			// Skip directories
			if info.IsDir() {
				return nil
			}

			// Delete if last modified more than 3 minutes ago
			if time.Since(info.ModTime()) > 3*time.Minute {
				err := os.Remove(path)
				if err != nil {
					log.Printf("Failed to remove %s: %v", path, err)
				} else {
					log.Printf("Deleted old temp file: %s", path)
				}
			}
			return nil
		})
		if err != nil {
			log.Printf("Error during cleanup: %v", err)
		}
	}
}

func StartCleanupTicker() {
	ticker := time.NewTicker(config.CleanupTime)
	defer ticker.Stop()

	for range ticker.C {
		log.Println("ticker envoked")
		CleanupFiles()
	}
}

// Cleanup old files (files older than 1 hour)
func CleanupFiles() {
	db := db.GetDB()

	log.Println("cleaning envoked")
	// Fetch records of files older than 1 hour
	rows, err := db.Query(`SELECT file_path FROM uploads WHERE last_download_at < ?`, time.Now().Add(-config.CleanupTime))
	if err != nil {
		log.Printf("Error fetching old files for cleanup: %v", err)
		return
	}
	defer rows.Close()	

	var filePath string
	var filePathsToDelete []string

	// Collect file paths to delete
	for rows.Next() {
		err := rows.Scan(&filePath)
		if err != nil {
			log.Printf("Error scanning file path: %v", err)
			continue
		}
		filePathsToDelete = append(filePathsToDelete, filePath)
	}

	// Delete files from the filesystem first
	for _, path := range filePathsToDelete {
		err := os.Remove(path)
		if err != nil {
			log.Printf("Error deleting file %s from filesystem: %v", path, err)
			continue
		}
		log.Printf("Successfully deleted file %s from filesystem.", path)
	}

	// After deleting from the filesystem, delete records from the database
	for _, path := range filePathsToDelete {
		_, err := db.Exec(`DELETE FROM uploads WHERE file_path = ?`, path)
		if err != nil {
			log.Printf("Error deleting file %s record from database: %v", path, err)
			continue
		}
		log.Printf("Successfully deleted file %s record from database.", path)
	}
}
