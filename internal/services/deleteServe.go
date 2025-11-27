package services

import (
	"fmt"
	"log"
	"os"
	"tiny-drop/internal/config"
	"tiny-drop/internal/db"
)



func DeleteUpload(fileUUID string) error {
    // Get the database connection
    conn := db.GetDB()
    table := config.UploadsTable

    // Look up the file object based on UUID
    fileObj, err := FileLookup(fileUUID)
    if err != nil {
        return fmt.Errorf("failed to lookup file with UUID %s: %w", fileUUID, err)
    }

    // Log the file path before checking existence
    log.Printf("Checking if file exists at path: %s", fileObj.FilePath)

    // Check if the file exists before trying to remove it
    if _, err := os.Stat(fileObj.FilePath); err == nil {
        // The file exists, so try to delete it
        if err := os.Remove(fileObj.FilePath); err != nil {
            return fmt.Errorf("failed to remove file %s: %w", fileObj.FilePath, err)
        }
        log.Printf("Successfully removed file %s", fileObj.FilePath)
    } else if os.IsNotExist(err) {
        // File doesn't exist, log the message
        log.Printf("file %s not found, skipping file removal", fileObj.FilePath)
    } else {
        // Other errors
        return fmt.Errorf("failed to check file existence for %s: %w", fileObj.FilePath, err)
    }

    // SQL query to delete the record from the database
    query := fmt.Sprintf(`DELETE FROM %s WHERE id = ?`, table)
	log.Printf("DELETE FROM %s WHERE AND id = %d", table, fileObj.ID)
    result, err := conn.Exec(query, fileObj.ID)
    if err != nil {
        return fmt.Errorf("failed to execute delete query for file %s: %w", fileUUID, err)
    }

    // Check how many rows were affected
    rows, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to get affected rows after deleting file %s: %w", fileUUID, err)
    }

    if rows != 1 {
        return fmt.Errorf("expected 1 row to be deleted, but got %d for file %s", rows, fileUUID)
    }

    // Return nil if everything was successful
    log.Printf("Successfully deleted file record with UUID %s from the database", fileUUID)
    return nil
}



func DeleteUpload_(fileUUID string) error {

	conn := db.GetDB()
	table := config.UploadsTable
	fileObj, err := FileLookup(fileUUID)
	if err != nil {
		return err
	}
	os.Remove(fileObj.FilePath)

	query := fmt.Sprintf(`DELETE FROM %s where file_uuid = ? and if = ?`, table)
	result, err := conn.Exec(query, fileUUID, fileObj.ID)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows != 1 {
		log.Printf("Affected row %d expected 1", rows)
	}

	return nil
}
