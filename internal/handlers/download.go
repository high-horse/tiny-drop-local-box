package handlers

import (
	"fmt"
	"net/http"
	"os"
	"tiny-drop/internal/services"
	"tiny-drop/internal/utils"
	"log"
	"io"
)


func HandleDownloadStream(w http.ResponseWriter, r *http.Request) {
	fileUUID := r.URL.Query().Get("fileUUID")
	if fileUUID == "" {
		utils.SendError(w, http.StatusUnprocessableEntity, "fileUUID is required.", nil)
		return
	}

	fileObj, err := services.FileLookup(fileUUID)
	if err != nil {
		utils.SendError(w, http.StatusNotFound, "File Not found.", nil)
		return
	}

	file, err := os.Open((*fileObj).FilePath)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Could not open file.", nil)
		return
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Cannot stat file.", nil)
		return
	}

	// Set headers
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", (*fileObj).FileName))
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", stat.Size()))

	// Stream file in chunks (e.g., 10MB)
	const chunkSize = 10 * 1024 * 1024
	buf := make([]byte, chunkSize)
	for {
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			utils.SendError(w, http.StatusInternalServerError, "Error reading file", nil)
			return
		}
		if n == 0 {
			break
		}

		if _, err := w.Write(buf[:n]); err != nil {
			log.Println("Client disconnected:", err)
			return
		}
		w.(http.Flusher).Flush()
	}
}




func HandleNormalDownload(w http.ResponseWriter, r *http.Request) {
	fileUUID := r.URL.Query().Get("fileUUID")
	if fileUUID == "" {
		utils.SendError(w, http.StatusUnprocessableEntity, "fileUUUID is required.", nil)
		return
	}

	fileObj, err := services.FileLookup(fileUUID)
	if err != nil {
		utils.SendError(w, http.StatusNotFound, "File Not found.", nil)
		return
	}

	file, err := os.Open((*fileObj).FilePath)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Could not open file.", nil)
		return
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Cannot stat file.", nil)
		return
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment;filename=%q", (*fileObj).FileName))

	http.ServeContent(w, r, stat.Name(), stat.ModTime(), file)
}

func HandleFileInfo(w http.ResponseWriter, r *http.Request) {
	fileUUID := r.URL.Query().Get("fileUUID")
	if fileUUID == "" {
		utils.SendError(w, http.StatusUnprocessableEntity, "fileUUID is required.", nil)
		return
	}

	fileObj, err := services.FileLookup(fileUUID)
	if err != nil {
		utils.SendError(w, http.StatusNotFound, "File Not found.", nil)
		return
	}

	utils.SendSuccess(w, http.StatusOK, "fetched successfully", (*fileObj))
}
