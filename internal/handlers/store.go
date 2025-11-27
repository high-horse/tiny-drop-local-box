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
	"time"
	"tiny-drop/internal/config"
	"tiny-drop/internal/db"
	"tiny-drop/internal/types"
	"tiny-drop/internal/utils"
)

const (
	tempDir  = config.TempDir
	finalDir = config.FinalDir
)


func UploadHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		utils.SendError(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
		// http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	err := r.ParseMultipartForm(10 << 20) // 10 mn buffer max
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "Failed to parse form", nil)
		return
		// http.Error(w, "Failed to parse form", http.StatusBadRequest)
	}
	file, _, err := r.FormFile("file")
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "Failed to Read File", nil)
		return
		// http.Error(w, "Failed to read file", http.StatusBadRequest)
	}
	defer file.Close()

	fileName := r.FormValue("fileName")
	chunkIndexStr := r.FormValue("chunkIndex")
	totalChunkStr := r.FormValue("totalChunks")
	totalSizeStr := r.FormValue("totalSize")

	chunkIndex, err := strconv.Atoi(chunkIndexStr)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "Invalid chunk index", nil)
		return
		// http.Error(w, "Invalid chunk index", http.StatusBadRequest)
	}
	totalChunks, err := strconv.Atoi(totalChunkStr)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "Invalid total chunks", nil)
		return
		// http.Error(w, "Invalid total chunks", http.StatusBadRequest)
	}

	totalSize, err := strconv.Atoi(totalSizeStr)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "Invalid total size", nil)
		return
		// http.Error(w, "Invalid total size", http.StatusBadRequest)
	}

	freeDiskSpace := utils.CheckDiskSpace(uint64(totalSize))
	log.Println("sisksopave ", freeDiskSpace)
	if chunkIndex == 0 && !freeDiskSpace {
		utils.SendError(w, http.StatusInsufficientStorage, "Server Space full, please try again later", nil)
		return
		// http.Error(w, "Server Space full, please try again later", http.StatusInternalServerError)
	}

	if chunkIndex < 0 || chunkIndex >= totalChunks {
		utils.SendError(w, http.StatusBadRequest, "Invalid chunk index.", nil)
		return
		// http.Error(w, "Invalid chunk index", http.StatusBadRequest)
	}

	// store in temp
	if err := os.MkdirAll(tempDir, os.ModePerm); err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Failed to create temp dir.", nil)
		return
		// http.Error(w, "Failed to create temp directory", http.StatusInternalServerError)
	}
	if err := os.MkdirAll(finalDir, os.ModePerm); err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Failed to create final Dir.", nil)
		return
		// http.Error(w, "Failed to create final directory", http.StatusInternalServerError)
	}

	// save the chunk
	uploadId := r.FormValue("uploadId")
	// chunkFilePath := filepath.Join(tempDir, fmt.Sprintf("%s.part%d", fileName, chunkIndex))
	chunkFilePath := filepath.Join(tempDir, fmt.Sprintf("%s.part%d", uploadId, chunkIndex))
	out, err := os.Create(chunkFilePath)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Failed to create chunk file.", nil)
		return
		// http.Error(w, "Failed to create chunk file", http.StatusInternalServerError)
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Failed to save chunk.", nil)
		return
		// http.Error(w, "Failed to save chunk", http.StatusInternalServerError)
	}

	if chunkIndex+1 == totalChunks {
		// finalPath := filepath.Join(finalDir, fileName)
		// finalFile, err := os.Create(finalPath)
		// if err != nil {
		// 	utils.SendError(w, http.StatusInternalServerError, "Failed to create final file.", nil)
		// 	return
		// 	// http.Error(w, "Failed to create final file", http.StatusInternalServerError)
		// }
		ext := filepath.Ext(fileName)
		finalFileName := uploadId + ext
		finalPath := filepath.Join(finalDir, finalFileName)
		finalFile, err := os.Create(finalPath)
		if err != nil {
			utils.SendError(w, http.StatusInternalServerError, "Failed to create final file.", nil)
			return
		}

		defer finalFile.Close()
		var fileSize int64

		for i := 0; i < totalChunks; i++ {
			// partPath := filepath.Join(tempDir, fmt.Sprintf("%s.part%d", fileName, i))
			partPath := filepath.Join(tempDir, fmt.Sprintf("%s.part%d", uploadId, i))
			partFile, err := os.Open(partPath)
			if err != nil {
				utils.SendError(w, http.StatusInternalServerError, "Failed to open chunk.", nil)
				return
				// http.Error(w, "Failed to open chunk", http.StatusInternalServerError)
			}
			n, err := io.Copy(finalFile, partFile)
			partFile.Close()
			if err != nil {
				utils.SendError(w, http.StatusInternalServerError, "failed to merge chunk.", nil)
				return
				// http.Error(w, "Failed to merge chunk", http.StatusInternalServerError)
			}
			fileSize += n
			os.Remove(partPath) // delete chunk after merging
		}

		log.Printf("File %s uploaded successfully", fileName)
		ip := utils.GetUserIp(r)
		
		metadata := types.FileMetadata{
			FileName:      fileName,
			UUID:          uploadId,
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

	utils.SendSuccess(w, http.StatusOK, "Saved successfully", nil)
	// w.WriteHeader(http.StatusOK)
	// w.Write([]byte("Chunk uploaded successfully"))
}

func SaveFileToDB(ip, fileUUID, filename, filepath string, filesize int64, metadataJSON []byte) error {
	db := db.GetDB()
	log.Println("Insertint to database ", ip, fileUUID, filename, filesize, metadataJSON)

	insertSQL := `
		INSERT INTO uploads 
		(ip, file_uuid, file_name, file_path, file_size, metadata, uploaded_at, last_download_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := db.Exec(insertSQL, ip, fileUUID, filename, filepath, filesize, metadataJSON, time.Now(), time.Now())
	if err != nil {
		return fmt.Errorf("failed to insert file into database: %v", err)
	}

	return nil
}
