package services

import (
	"database/sql"
	"fmt"
	"log"
	"tiny-drop/internal/config"
	"tiny-drop/internal/db"
	"tiny-drop/internal/types"
)

func FileLookup(fileUUID string) (*types.UploadFile, error) {
	conn := db.GetDB()
	table := config.UploadsTable

	// Log the query to ensure it's being formatted correctly
	log.Printf("SELECT id, ip, file_uuid, file_name, file_path, file_size, uploaded_at, last_download_at, uploader_id FROM %s WHERE file_uuid = ?\n", table)
	query := fmt.Sprintf(`
		SELECT ip, file_uuid, file_name, file_path, file_size, uploaded_at, last_download_at, uploader_id
		FROM %s
		WHERE file_uuid = ?`, table)

	var uploadFile types.UploadFile

	// Pass the fileUUID as a parameter to the query (avoiding SQL injection risks)
	row := conn.QueryRow(query, fileUUID)
	err := row.Scan(
		// &uploadFile.ID,
		&uploadFile.IP,
		&uploadFile.FileUUID,
		&uploadFile.FileName,
		&uploadFile.FilePath,
		&uploadFile.FileSize,
		&uploadFile.UploadedAt,
		&uploadFile.LastDownloadAt,
		&uploadFile.UploaderId,
	)

	if err != nil {
		log.Printf("Error during Scan: %v", err)
		if err == sql.ErrNoRows {
			log.Printf("No file found with UUID: %s", fileUUID)
			return nil, fmt.Errorf("no file found with UUID: %s", fileUUID)
		}
		return nil, fmt.Errorf("error scanning row: %v", err)
	}

	// Log the successfully fetched file data
	log.Printf("Fetched file: %+v", uploadFile)

	return &uploadFile, nil
}




func FileLookup_(fileUUID string) (*types.UploadFile, error) {
	conn := db.GetDB()
	table := config.UploadsTable

	// Use placeholder for the table name (which is safe)
	log.Printf("SELECT id, ip, file_uuid, file_name, file_path, file_size, uploaded_at, last_download_at, metadata, uploader_id FROM %s WHERE file_uuid = ?\n\n\n", table)
	query := fmt.Sprintf(`
		SELECT id, ip, file_uuid, file_name, file_path, file_size, uploaded_at, last_download_at, metadata, uploader_id
		FROM %s
		WHERE file_uuid = ?`, table)

	var uploadFile types.UploadFile

	// Pass the fileUUID as a parameter to the query (avoiding SQL injection risks)
	row := conn.QueryRow(query, fileUUID)
	err := row.Scan(
		&uploadFile.ID,
		&uploadFile.IP,
		&uploadFile.FileUUID,
		&uploadFile.FileName,
		&uploadFile.FilePath,
		&uploadFile.FileSize,
		&uploadFile.UploadedAt,
		&uploadFile.LastDownloadAt,
		&uploadFile.Metadata,  // Make sure your metadata is properly deserialized from JSON
		&uploadFile.UploaderId,
		&uploadFile.WillDeleteAt,
	)
	if err != nil {
		return nil, fmt.Errorf("no file found with UUID: %s", fileUUID)
	}

	return &uploadFile, nil
}
