package main

import (
	"database/sql"
	"file-management/handlers"

	// "file-management/jobs" // Import the jobs package
	"file-management/jobs"
	"file-management/middleware"
	"file-management/utils"
	"log"
	"net/http"
	"os"

	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var db *sql.DB
var jwtSecret string

func main() {
	utils.InitRedis()

	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Retrieve JWT secret from environment variable
	jwtSecret = os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable not set")
	}

	connStr := "user=" + os.Getenv("POSTGRES_USER") + " password=" + os.Getenv("POSTGRES_PASSWORD") + " dbname=file_sharing sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}
	log.Println("Database connection established")

	// Initialize the background job
	jobs.Init(connStr)
	go jobs.RunBackgroundJob(1 * time.Minute) // Run the job every 1 minute for testing

	// Router configuration
	r := mux.NewRouter()
	r.HandleFunc("/register", handlers.RegisterHandler).Methods("POST")
	r.HandleFunc("/login", handlers.LoginHandler).Methods("POST")
	r.Handle("/upload", middleware.AuthMiddleware(http.HandlerFunc(handlers.UploadFileHandler)))
	r.Handle("/delete", middleware.AuthMiddleware(http.HandlerFunc(handlers.DeleteFileHandler)))
	r.Handle("/files/{id:[0-9]+}", middleware.AuthMiddleware(http.HandlerFunc(handlers.UpdateFileHandler))).Methods("PUT")
	r.Handle("/files", middleware.AuthMiddleware(http.HandlerFunc(handlers.SearchFilesHandler))).Methods("GET")
	r.Handle("/share", middleware.AuthMiddleware(http.HandlerFunc(handlers.ShareFileHandler)))

	log.Fatal(http.ListenAndServe(":8080", r))
}
