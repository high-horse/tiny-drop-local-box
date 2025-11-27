package services

import (
	"encoding/json"
	"fmt"
	"log"
	"tiny-drop/internal/config"
	"tiny-drop/internal/db"
	"tiny-drop/internal/types"
)


func FetchUploadeds(ip string) ([]types.UploadFile, error) {
	conn := db.GetDB()
	table := config.UploadsTable

	query := fmt.Sprintf(`
		SELECT id, ip, file_uuid, file_name, file_path, file_size, uploaded_at, last_download_at, metadata
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
		)
		if err != nil {
			log.Printf("Failed to scan row: %v", err)
			continue
		}

		uf.FileSize = fmt.Sprintf("%d", fileSize)

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