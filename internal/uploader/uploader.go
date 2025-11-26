package uploader

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"
	"tiny-drop/internal/db"
)

type FileMetadata struct {
	FileType      string   `json:"file_type"`
	FileSize      int64    `json:"file_size"`
	FileExtension string   `json:"file_extension"`
	UploadedBy    string   `json:"uploaded_by,omitempty"`
	Tags          []string `json:"tags,omitempty"`
}

func UploadFile(w http.ResponseWriter, r *http.Request) {
	ip := GetUserIp(r)

	err := r.ParseMultipartForm(10 << 20) // 10MB buffer size
	if err != nil {
		http.Error(w, "Error parsing the file", http.StatusBadRequest)
		return
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error reading the file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	fileInfo, err := fileHeader.Open()
	if err != nil {
		http.Error(w, "Error opening the file", http.StatusInternalServerError)
		return
	}
	defer fileInfo.Close()

	// File size in bytes
	fileSize := fileHeader.Size
	if !CheckDiskSpace(fileSize) {
		http.Error(w, "Not enough disk space to upload the file.", http.StatusInternalServerError)
		return
	}

	fileUUID := uuid.New().String()
	ext := filepath.Ext(fileHeader.Filename)
	filename := fmt.Sprintf("%s%s", fileUUID, ext)
	filePath := filepath.Join("./storage/uploads", filename)

	outFile, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Error saving the file", http.StatusInternalServerError)
		return
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, file)
	if err != nil {
		http.Error(w, "Error saving file to disk", http.StatusInternalServerError)
		return
	}

	metadata := FileMetadata{
		FileType:      fileHeader.Header.Get("Content-Type"),
		FileSize:      fileSize,
		FileExtension: ext,
		UploadedBy:    ip,                        // Could be the IP address or any other info
		Tags:          []string{"user-uploaded"}, // Example tag, you can modify as needed
	}

	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		http.Error(w, "Error serializing metadata to JSON", http.StatusInternalServerError)
		return
	}

	err = SaveFileToDB(ip, fileUUID, filename, filePath, fileSize, metadataJSON)

	if err != nil {
		http.Error(w, "Error saving file info to database", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "File uploaded successfully: %s", filename)

	// go cleanupFiles()
}

func GetUserIp(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	var err error
	if ip == "" {
		ip = r.Header.Get("X-Real-IP")
	}
	if ip == "" {
		ip, _, err = net.SplitHostPort(ip)
		if err != nil {
			log.Printf("Error Parsing IP address : %v", err)
			return ""
		}
	}

	return ip
}

// Get the user IP address from the request (considering proxy headers)
func getUserIP(r *http.Request) string {

	/*
		NGINX
		proxy_set_header X-Forwarded-For $remote_addr;
		proxy_set_header X-Real-IP $remote_addr;
		proxy_set_header Host $host;
		proxy_pass http://your_go_app_backend;
	*/

	// Check if the request has the X-Forwarded-For header
	// The X-Forwarded-For header contains a comma-separated list of IPs
	// The first IP in the list is usually the original client IP
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		// The first IP in the list is the client's IP
		ips := strings.Split(xff, ",")
		return strings.TrimSpace(ips[0])
	}

	// If no X-Forwarded-For header, use RemoteAddr as fallback
	ip := r.RemoteAddr
	host, _, err := net.SplitHostPort(ip)
	if err != nil {
		log.Println("Error parsing IP address:", err)
		return ""
	}
	return host
}

func CheckDiskSpace(fileSize int64) bool {
	statfs := syscall.Statfs_t{}
	err := syscall.Statfs("./storage", &statfs)
	if err != nil {
		log.Println("Error checking disk space:", err)
		return false
	}

	// Get the available free space (in bytes)
	freeSpace := statfs.Bavail * uint64(statfs.Bsize)

	// Check if the available space after the upload will be enough
	if freeSpace < uint64(fileSize) {
		log.Printf("Not enough disk space! Required: %d bytes, Available: %d bytes", fileSize, freeSpace)
		return false
	}

	// If there's enough space
	return true
}

func SaveFileToDB(ip, fileUUID, filename, filepath string, filesize int64, metadataJSON []byte) error {
	db := db.GetDB()

	insertSQL := `
		INSERT INTO uploads 
		(ip, file_uuid, filename, file_path, filesize, uploaded_at, metadata) 
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	_, err := db.Exec(insertSQL, ip, fileUUID, filename, filepath, filesize, time.Now(), metadataJSON)
	if err != nil {
		return fmt.Errorf("failed to insert file into database: %v", err)
	}

	return nil
}

func SaveFileToDB_(ip, filename, filepath string, filesize int64, metadataJSON []byte) error {
	db := db.GetDB()

	insertSql := `
		INSERT INTO uploads
		(ip, filename, filepath, filesize, metadata)
		VALUES
		(?, ?, ?, ?, ?)
	`

	_, err := db.Exec(insertSql, ip, filename, filepath, filesize, metadataJSON)
	if err != nil {
		return fmt.Errorf("failed to insert file to database : %v", err)
	}
	return nil
}
