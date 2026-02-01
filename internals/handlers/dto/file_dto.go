package dto

type FileResponse struct {
	DocID string `json:"doc_id"`
}

type DocRequest struct {
	FileName      string `json:"file_name"`
	MiMeType      string `json:"mime_type"`
	FileSize      int64  `json:"file_size"`
	UploadPurpose string `json:"upload_purpose"` // may be enum
}

type PresignedResponse struct {
	DocID     string `json:"doc_id"`
	UploadURL string `json:"upload_url"`
}
