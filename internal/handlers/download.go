package handlers

import (
	"fmt"
	"net/http"
	"os"
	"tiny-drop/internal/services"
	"tiny-drop/internal/utils"
)


func HandleDownloadStream(w http.ResponseWriter, r *http.Request) {
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
