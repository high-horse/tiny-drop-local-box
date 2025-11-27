package handlers

import (
	"log"
	"net/http"
	"tiny-drop/internal/services"
	"tiny-drop/internal/utils"
)


func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	fileUUID := r.URL.Query().Get("fileUUID")
	if fileUUID == "" {
		utils.SendError(w, http.StatusUnprocessableEntity, "fileUUID is required.", nil)
		return
	}

	if err := services.DeleteUpload(fileUUID); err != nil {
		log.Println("Failed to delete file with uuid ")
		utils.SendError(w, http.StatusNotFound, "File Not found.", nil)
		return
	}

	utils.SendSuccess(w, http.StatusOK, "Successfully deleted", nil)

}