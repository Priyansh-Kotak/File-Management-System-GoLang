package handlers

import (
	"file-management/utils"
	"io"
	"net/http"
	"os"
)

// UploadFileHandler handles file uploads
func UploadFileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Extract the JWT token from the Authorization header
	tokenString := r.Header.Get("Authorization")
	_, err := utils.VerifyJWT(tokenString)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// File handling logic (save locally)
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error reading file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileData, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Error reading file data", http.StatusInternalServerError)
		return
	}

	// Save file to the uploads directory
	fileName := "example.txt" // In production, generate a unique name
	err = os.WriteFile("./uploads/"+fileName, fileData, 0644)
	if err != nil {
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("File uploaded successfully"))
}
