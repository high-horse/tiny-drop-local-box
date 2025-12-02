package handlers

import (
	"encoding/json"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"tiny-drop/internal/types"
	"tiny-drop/internal/utils"
)


func UploadStreamHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.SendError(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
		// http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
	uploadID := r.Header.Get("X-Upload-ID")
	fileName := r.Header.Get("X-File-Name")
	uploaderID := r.Header.Get("X-Uploader-ID")

	if uploadID == "" || fileName == "" {
		utils.SendError(w, http.StatusBadRequest, "Missing upload headers", nil)
		return
	}

	ext := filepath.Ext(fileName)
	finalFileName := uploadID + ext

	// Ensure temp and final directories exist
	if err := os.MkdirAll(tempDir, os.ModePerm); err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Failed to create temp dir", nil)
		return
	}
	if err := os.MkdirAll(finalDir, os.ModePerm); err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Failed to create final dir", nil)
		return
	}

	// Write to temp file first
	tempFilePath := filepath.Join(tempDir, finalFileName)
	tempFile, err := os.Create(tempFilePath)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Failed to create temp file", nil)
		return
	}
	defer tempFile.Close()

	// Stream body into temp file
	n, err := io.Copy(tempFile, r.Body)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Failed to write file", nil)
		return
	}

	// Move temp file to final directory
	finalPath := filepath.Join(finalDir, finalFileName)
	if err := os.Rename(tempFilePath, finalPath); err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Failed to move file to final directory", nil)
		return
	}

	ip := utils.GetUserIp(r)
	metadata := types.FileMetadata{
		FileName:      fileName,
		UUID:          uploadID,
		FileType:      mime.TypeByExtension(ext),
		FileSize:      n,
		FileExtension: ext,
		UploadedBy:    ip,
		Tags:          []string{"user-uploaded"},
	}

	metadataJSON, _ := json.Marshal(metadata)
	SaveFileToDB(ip, uploadID, fileName, finalPath, n, metadataJSON, uploaderID)

	log.Println("Stream upload completed:", finalPath)
	utils.SendSuccess(w, http.StatusOK, "File streamed successfully", nil)
}

func UploadStreamHandlerFinalDir(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		utils.SendError(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}

	uploadID := r.Header.Get("X-Upload-ID")
    fileName := r.Header.Get("X-File-Name")

	if uploadID == "" || fileName == "" {
        utils.SendError(w, http.StatusBadRequest, "Missing upload headers", nil)
        return
    }

	ext := filepath.Ext(fileName)
    finalFileName := uploadID + ext

	
	finalPath := filepath.Join(finalDir, finalFileName)
    finalFile, err := os.Create(finalPath)
    if err != nil {
        utils.SendError(w, 500, "Failed to create file", nil)
        return
    }
    defer finalFile.Close()

	// Stream input directly into file
    n, err := io.Copy(finalFile, r.Body)
    if err != nil {
        utils.SendError(w, 500, "Failed to write file", nil)
        return
    }

	ip := utils.GetUserIp(r)
	metadata := types.FileMetadata{
        FileName:      fileName,
        UUID:          uploadID,
        FileType:      mime.TypeByExtension(ext),
        FileSize:      n,
        FileExtension: ext,
        UploadedBy:    ip,
        Tags:          []string{"user-uploaded"},
    }

    metadataJSON, _ := json.Marshal(metadata)

	SaveFileToDB(ip, uploadID, fileName, finalPath, n, metadataJSON, "") // uploaderId optional

    log.Println("Stream upload completed:", finalPath)

    utils.SendSuccess(w, http.StatusOK, "File streamed successfully", nil)
}
