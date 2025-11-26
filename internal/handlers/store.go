package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"tiny-drop/internal/config"
	"tiny-drop/internal/db"
	"tiny-drop/internal/utils"
)

const (
	tempDir  = config.TempDir
	finalDir = config.FinalDir
)

type FileMetadata struct {
	FileType      string   `json:"file_type"`
	FileSize      int64    `json:"file_size"`
	FileExtension string   `json:"file_extension"`
	UploadedBy    string   `json:"uploaded_by,omitempty"`
	Tags          []string `json:"tags,omitempty"`
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(10 << 20) // 10 mn buffer max
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to read file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileName := r.FormValue("fileName")
	chunkIndexStr := r.FormValue("chunkIndex")
	totalChunkStr := r.FormValue("totalChunks")

	chunkIndex, err := strconv.Atoi(chunkIndexStr)
	if err != nil {
		http.Error(w, "Invalid chunk index", http.StatusBadRequest)
		return
	}
	totalChunks, err := strconv.Atoi(totalChunkStr)
	if err != nil {
		http.Error(w, "Invalid total chunks", http.StatusBadRequest)
		return
	}
	if chunkIndex < 0 || chunkIndex >= totalChunks {
		http.Error(w, "Invalid chunk index", http.StatusBadRequest)
		return
	}

	// store in temp
	if err := os.MkdirAll(tempDir, os.ModePerm); err != nil {
		http.Error(w, "Failed to create temp directory", http.StatusInternalServerError)
		return
	}
	if err := os.MkdirAll(finalDir, os.ModePerm); err != nil {
		http.Error(w, "Failed to create final directory", http.StatusInternalServerError)
		return
	}

	// save the chunk
	uploadId := r.FormValue("uploadId")
	// chunkFilePath := filepath.Join(tempDir, fmt.Sprintf("%s.part%d", fileName, chunkIndex))
	chunkFilePath := filepath.Join(tempDir, fmt.Sprintf("%s.part%d", uploadId, chunkIndex))
	out, err := os.Create(chunkFilePath)
	if err != nil {
		http.Error(w, "Failed to create chunk file", http.StatusInternalServerError)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		http.Error(w, "Failed to save chunk", http.StatusInternalServerError)
		return
	}

	if chunkIndex+1 == totalChunks {
		finalPath := filepath.Join(finalDir, fileName)
		finalFile, err := os.Create(finalPath)
		if err != nil {
			http.Error(w, "Failed to create final file", http.StatusInternalServerError)
			return
		}
		defer finalFile.Close()
		var fileSize int64

		for i := 0; i < totalChunks; i++ {
			// partPath := filepath.Join(tempDir, fmt.Sprintf("%s.part%d", fileName, i))
			partPath := filepath.Join(tempDir, fmt.Sprintf("%s.part%d", uploadId, i))
			partFile, err := os.Open(partPath)
			if err != nil {
				http.Error(w, "Failed to open chunk", http.StatusInternalServerError)
				return
			}
			n, err := io.Copy(finalFile, partFile)
			partFile.Close()
			if err != nil {
				http.Error(w, "Failed to merge chunk", http.StatusInternalServerError)
				return
			}
			fileSize += n
			os.Remove(partPath) // delete chunk after merging
		}

		log.Printf("File %s uploaded successfully", fileName)
		ip := utils.GetUserIp(r)
		ext := filepath.Ext(fileName)

		metadata := FileMetadata{
			FileType:      mime.TypeByExtension(ext),
			FileSize:      fileSize,
			FileExtension: ext,
			UploadedBy:    ip,
			Tags:          []string{"user-uploaded"},
		}
		metadataJSON, err := json.Marshal(metadata)
		if err != nil {
			log.Printf("Failed to marshal file metadata: %v", err)
			metadataJSON = []byte("{}") // fallback empty JSON
		}
		SaveFileToDB(ip, uploadId, fileName, finalPath, fileSize, metadataJSON)
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Chunk uploaded successfully"))
}

func SaveFileToDB(ip, fileUUID, filename, filepath string, filesize int64, metadataJSON []byte) error {
	db := db.GetDB()

	insertSQL := `
		INSERT INTO uploads 
		(ip, file_uuid, file_name, file_path, file_size, metadata) 
		VALUES (?, ?, ?, ?, ?, ?)gi
	`
	_, err := db.Exec(insertSQL, ip, fileUUID, filename, filepath, filesize, metadataJSON)
	if err != nil {
		return fmt.Errorf("failed to insert file into database: %v", err)
	}

	return nil
}
