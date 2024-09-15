package models

type File struct {
    UserID   int    `json:"user_id"`
    FileName string `json:"file_name"`
}