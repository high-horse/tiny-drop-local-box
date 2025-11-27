package types

import "time"

type FileMetadata struct {
	FileName      string   `json:"file_name"`
	UUID          string   `json:"uuid"`
	FileType      string   `json:"file_type"`
	FileSize      int64    `json:"file_size"`
	FileExtension string   `json:"file_extension"`
	UploadedBy    string   `json:"uploaded_by,omitempty"`
	Tags          []string `json:"tags,omitempty"`
}

type UploadFile struct {
	ID             int          `json:"id"`
	IP             string       `json:"ip"`
	FileUUID       string       `json:"file_uuid"`
	FileName       string       `json:"file_name"`
	FilePath       string       `json:"file_path"`
	FileSize       string       `json:"file_size"`
	UploadedAt     time.Time    `json:"uploaded_at"`
	LastDownloadAt time.Time    `json:"last_download_at"`
	Metadata       FileMetadata `json:"metadata"`
}
