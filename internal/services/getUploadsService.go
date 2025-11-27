package services

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
	"tiny-drop/internal/config"
	"tiny-drop/internal/db"
	"tiny-drop/internal/types"
)


func FetchUploadeds(ip string) ([]types.UploadFile, error) {
	conn := db.GetDB()
	table := config.UploadsTable

	query := fmt.Sprintf(`
		SELECT id, ip, file_uuid, file_name, file_path, file_size, uploaded_at, last_download_at, metadata, uploader_id
		FROM %s
		WHERE ip = ?`, table,
	)

	rows, err := conn.Query(query, ip)
	if err != nil {
		return nil, fmt.Errorf("failed to query uploads: %v", err)
	}
	defer rows.Close()

	var uploads []types.UploadFile

	for rows.Next() {
		var uf types.UploadFile
		var metadataJSON []byte
		var fileSize int64

		err := rows.Scan(
			&uf.ID,
			&uf.IP,
			&uf.FileUUID,
			&uf.FileName,
			&uf.FilePath,
			&fileSize,
			&uf.UploadedAt,
			&uf.LastDownloadAt,
			&metadataJSON,
			&uf.UploaderId,
		)
		if err != nil {
			log.Printf("Failed to scan row: %v", err)
			continue
		}

		uf.FileSize = formatFileSize(fileSize)
		uf.WillDeleteAt = uf.LastDownloadAt.Add(1 * time.Hour)

		// Unmarshal metadata JSON
		var metadata types.FileMetadata
		if err := json.Unmarshal(metadataJSON, &metadata); err != nil {
			log.Printf("Failed to unmarshal metadata for file %s: %v", uf.FileUUID, err)
			metadata = types.FileMetadata{}
		}
		uf.Metadata = metadata

		uploads = append(uploads, uf)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	return uploads, nil
}


func formatFileSize(fileSize int64) string {
	var size float64
	var unit string

	// Determine the unit and convert the size accordingly
	if fileSize < 1000 {
		size = float64(fileSize)
		unit = "B" // Bytes
	} else if fileSize < 1000000 {
		size = float64(fileSize) / 1000
		unit = "KB" // Kilobytes
	} else if fileSize < 1000000000 {
		size = float64(fileSize) / 1000000
		unit = "MB" // Megabytes
	} else {
		size = float64(fileSize) / 1000000000
		unit = "GB" // Gigabytes
	}

	// Format the size to 2 decimal places and append the unit
	return fmt.Sprintf("%.2f %s", size, unit)
}