package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	_ "github.com/lib/pq" // Importing the driver
)

var db *sql.DB // Initialize your database connection

type File struct {
	ID       int    `json:"id"`
	UserID   int    `json:"user_id"`
	FileName string `json:"file_name"`
	FileSize int    `json:"file_size"`
	S3URL    string `json:"s3_url"`
}

func UploadFileHandler(w http.ResponseWriter, r *http.Request) {
	userEmail, ok := r.Context().Value("userEmail").(string)
	if !ok || userEmail == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Retrieve user ID from the database
	var userID int
	err := db.QueryRow("SELECT id FROM users WHERE email = $1", userEmail).Scan(&userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	// Retrieve file metadata from the request
	var file File
	err = json.NewDecoder(r.Body).Decode(&file)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Check that required fields are present
	if file.FileName == "" || file.FileSize <= 0 || file.S3URL == "" {
		http.Error(w, "Missing file metadata", http.StatusBadRequest)
		return
	}

	// Save file metadata in the database with userID
	_, err = db.Exec("INSERT INTO files (user_id, file_name, file_size, s3_url) VALUES ($1, $2, $3, $4)", userID, file.FileName, file.FileSize, file.S3URL)
	if err != nil {
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("File uploaded successfully"))
}

func DeleteFileHandler(w http.ResponseWriter, r *http.Request) {
	userEmail := r.Context().Value("userEmail").(string)

	// Retrieve user ID from the database
	var userID int
	err := db.QueryRow("SELECT id FROM users WHERE email = $1", userEmail).Scan(&userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	// Retrieve file ID from the URL
	fileID := r.URL.Query().Get("file_id")
	if fileID == "" {
		http.Error(w, "Missing file ID", http.StatusBadRequest)
		return
	}

	// Check if the file belongs to the user
	var fileOwnerID int
	err = db.QueryRow("SELECT user_id FROM files WHERE id = $1", fileID).Scan(&fileOwnerID)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	if fileOwnerID != userID {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Perform file deletion logic
	_, err = db.Exec("DELETE FROM files WHERE id = $1", fileID)
	if err != nil {
		http.Error(w, "Error deleting file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("File deleted successfully"))
}
