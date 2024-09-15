package handlers

import (
	"database/sql"
	"encoding/json"
	"file-management/cache"
	"file-management/models"
	"file-management/utils"
	"fmt"
	"io" // Import utility functions for file handling
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"github.com/go-redis/redis/v8"
)

var db *sql.DB // Initialize your database connection

type File struct {
	ID         int            `json:"id"`
	FileName   string         `json:"file_name"`
	FileSize   int            `json:"file_size"`
	S3URL      string         `json:"s3_url"`
	UploadDate time.Time      `json:"upload_date"`
	FileType   sql.NullString `json:"file_type"`
}

func UploadFileHandler(w http.ResponseWriter, r *http.Request) {
	userEmail := r.Context().Value("userEmail").(string)

	// Retrieve user ID from the database
	var userID int
	err := db.QueryRow("SELECT id FROM users WHERE email = $1", userEmail).Scan(&userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	// Parse the multipart form to retrieve the uploaded file
	err = r.ParseMultipartForm(10 << 20) // Max file size: 10MB
	if err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Create a directory to store the file if it doesn't exist
	uploadDir := "./upload"
	err = os.MkdirAll(uploadDir, os.ModePerm)
	if err != nil {
		http.Error(w, "Unable to create upload directory", http.StatusInternalServerError)
		return
	}

	// Create a new file in the upload directory
	dst, err := os.Create(uploadDir + "/" + handler.Filename)
	if err != nil {
		http.Error(w, "Unable to create the file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copy the uploaded file's content to the newly created file
	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, "Error saving the file", http.StatusInternalServerError)
		return
	}

	// Save file metadata in the database
	_, err = db.Exec("INSERT INTO files (user_id, file_name, file_size, s3_url) VALUES ($1, $2, $3, $4)",
		userID, handler.Filename, handler.Size, "/upload/"+handler.Filename)
	if err != nil {
		http.Error(w, "Error saving file metadata", http.StatusInternalServerError)
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

var ctx = utils.Ctx

func GetFilesHandler(w http.ResponseWriter, r *http.Request) {

	userEmail := r.Context().Value("userEmail").(string)

	// Retrieve user ID from database
	var userID int
	err := db.QueryRow("SELECT id FROM users WHERE email = $1", userEmail).Scan(&userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	// Check Redis cache first
	cacheKey := "files_user_" + strconv.Itoa(userID)
	cachedData, err := utils.RedisClient.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		// If not in cache, query the database
		rows, err := db.Query("SELECT id, file_name, file_size, s3_url FROM files WHERE user_id = $1", userID)
		if err != nil {
			http.Error(w, "Error retrieving files", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var files []File
		for rows.Next() {
			var file File
			if err := rows.Scan(&file.ID, &file.FileName, &file.FileSize, &file.S3URL); err != nil {
				http.Error(w, "Error scanning file data", http.StatusInternalServerError)
				return
			}
			files = append(files, file)
		}

		// Cache metadata in Redis
		fileData, _ := json.Marshal(files)
		utils.RedisClient.Set(ctx, cacheKey, fileData, 0)

		// Send response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(files)
	} else if err != nil {
		http.Error(w, "Error fetching cache", http.StatusInternalServerError)
	} else {
		// Return cached data
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(cachedData))
	}
}

func ShareFileHandler(w http.ResponseWriter, r *http.Request) {
	// Get file ID from URL path
	fileID := r.URL.Query().Get("file_id")
	if fileID == "" {
		http.Error(w, "Missing file ID", http.StatusBadRequest)
		return
	}

	// Retrieve file metadata from the database
	var file File
	err := db.QueryRow("SELECT file_name, s3_url FROM files WHERE id = $1", fileID).Scan(&file.FileName, &file.S3URL)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	// Send the public link for sharing
	publicURL := file.S3URL
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(publicURL))
}

var redisClient = cache.GetClient()

func SearchFilesHandler(w http.ResponseWriter, r *http.Request) {
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
		fmt.Println("Error retrieving user ID:", err)
		return
	}

	// Retrieve search parameters from query
	fileName := r.URL.Query().Get("file_name")
	uploadDate := r.URL.Query().Get("upload_date")
	fileType := r.URL.Query().Get("file_type")

	// Build the query
	query := "SELECT id, file_name, file_size, s3_url, upload_date, file_type FROM files WHERE user_id = $1"
	args := []interface{}{userID}

	if fileName != "" {
		query += " AND file_name ILIKE $" + strconv.Itoa(len(args)+1)
		args = append(args, "%"+fileName+"%")
	}
	if uploadDate != "" {
		query += " AND upload_date::date = $" + strconv.Itoa(len(args)+1)
		args = append(args, uploadDate)
	}
	if fileType != "" {
		query += " AND file_type = $" + strconv.Itoa(len(args)+1)
		args = append(args, fileType)
	}

	fmt.Println("Query:", query)
	fmt.Println("Args:", args)

	rows, err := db.Query(query, args...)
	if err != nil {
		http.Error(w, "Error retrieving files", http.StatusInternalServerError)
		fmt.Println("Error executing query:", err)
		return
	}
	defer rows.Close()

	var files []File
	for rows.Next() {
		var file File
		if err := rows.Scan(&file.ID, &file.FileName, &file.FileSize, &file.S3URL, &file.UploadDate, &file.FileType); err != nil {
			http.Error(w, "Error scanning files", http.StatusInternalServerError)
			fmt.Println("Error scanning rows:", err)
			return
		}
		files = append(files, file)
	}

	if len(files) == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode([]File{})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(files)
}

func UpdateFileHandler(w http.ResponseWriter, r *http.Request) {
	// Extract file ID from URL parameters
	vars := mux.Vars(r)
	fileID := vars["id"]

	// Parse the request body
	var req models.UpdateFileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Update file metadata in the database
	_, err := db.ExecContext(ctx, "UPDATE files SET file_name = $1 WHERE id = $2", req.FileName, fileID)
	if err != nil {
		http.Error(w, "Error updating file", http.StatusInternalServerError)
		return
	}

	// Invalidate cache
	cacheKey := "file:" + fileID
	err = redisClient.Del(ctx, cacheKey).Err()
	if err != nil {
		log.Println("Error invalidating cache:", err)
	}

	w.WriteHeader(http.StatusNoContent)
}
