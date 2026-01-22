package dto

type FileResponse struct {
	DocID string `json:"doc_id"`
}

type DocRequest struct {
	FileName string `json:"file_name"`
	MiMeType string `json:"mime_type"`
	FileSize int16 `json:"file_size"`
	UploadPurpose string `json:"upload_purpose"` // may be enum
}

type PresignedResponse struct{
	UploadURL string `json:"upload_url"`
}