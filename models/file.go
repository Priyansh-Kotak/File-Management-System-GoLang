package models

type File struct {
	UserID   int    `json:"user_id"`
	FileName string `json:"file_name"`
}

type UpdateFileRequest struct {
	FileName   *string `json:"file_name,omitempty"`
	FileSize   *int    `json:"file_size,omitempty"`
	S3URL      *string `json:"s3_url,omitempty"`
	UploadDate *string `json:"upload_date,omitempty"`
	FileType   *string `json:"file_type,omitempty"`
	UserID     int     `json:"user_id"`
	ID         int     `json:"id"`
 
}

