// // package main

// // import (
// // 	"database/sql"
// // 	"file-management/handlers"
// // 	"log"
// // 	"net/http"
// // 	"os"

// // 	"github.com/joho/godotenv"
// // 	_ "github.com/lib/pq"
// // 	// "github.com/yourusername/yourproject/handlers"
// // )

// // var db *sql.DB
// // var jwtSecret string

// // func main() {
// // 	// Load environment variables from .env file
// // 	err := godotenv.Load()
// // 	if err != nil {
// // 		log.Fatalf("Error loading .env file: %v", err)
// // 	}

// // 	// Retrieve JWT secret from environment variable
// // 	jwtSecret = os.Getenv("JWT_SECRET")
// // 	if jwtSecret == "" {
// // 		log.Fatal("JWT_SECRET environment variable not set")
// // 	}

// // 	// Update this connection string with your credentials
// // 	connStr := "user=" + os.Getenv("POSTGRES_USER") + " password=" + os.Getenv("POSTGRES_PASSWORD") + " dbname=file_sharing sslmode=disable"
// // 	db, err = sql.Open("postgres", connStr)
// // 	if err != nil {
// // 		log.Fatalf("Error opening database: %v", err)
// // 	}

// // 	// Test the connection
// // 	err = db.Ping()
// // 	if err != nil {
// // 		log.Fatalf("Error pinging database: %v", err)
// // 	}
// // 	log.Println("Database connection established")

// // 	http.HandleFunc("/register", handlers.RegisterHandler)
// // 	http.ListenAndServe(":8080", nil)
// // }

// package main

// import (
// 	"database/sql"
// 	"file-management/handlers"
// 	"file-management/middleware"
// 	"log"
// 	"net/http"
// 	"os"

// 	"github.com/joho/godotenv"
// 	_ "github.com/lib/pq"
// 	// "github.com/yourusername/yourproject/handlers"
// )

// var db *sql.DB
// var jwtSecret string

// func main() {
// 	// Load environment variables from .env file
// 	err := godotenv.Load()
// 	if err != nil {
// 		log.Fatalf("Error loading .env file: %v", err)
// 	}

// 	// Retrieve JWT secret from environment variable
// 	jwtSecret = os.Getenv("JWT_SECRET")
// 	if jwtSecret == "" {
// 		log.Fatal("JWT_SECRET environment variable not set")
// 	}

// 	connStr := "user=" + os.Getenv("POSTGRES_USER") + " password=" + os.Getenv("POSTGRES_PASSWORD") + " dbname=file_sharing sslmode=disable"
// 	db, err = sql.Open("postgres", connStr)
// 	if err != nil {
// 		log.Fatalf("Error opening database: %v", err)
// 	}

// 	// Test the connection
// 	err = db.Ping()
// 	if err != nil {
// 		log.Fatalf("Error pinging database: %v", err)
// 	}
// 	log.Println("Database connection established")

// 	// http.HandleFunc("/register", handlers.RegisterHandler)
// 	// http.HandleFunc("/login", handlers.LoginHandler)
// 	// http.Handle("/upload", middleware.AuthMiddleware(http.HandlerFunc(handlers.UploadFileHandler)))
// 	// http.Handle("/delete", middleware.AuthMiddleware(http.HandlerFunc(handlers.DeleteFileHandler)))
// 	http.Handle("/register", http.HandlerFunc(handlers.RegisterHandler))
// 	http.Handle("/login", http.HandlerFunc(handlers.LoginHandler))

// 	// Use the AuthMiddleware for file management routes
// 	http.Handle("/upload", middleware.AuthMiddleware(http.HandlerFunc(handlers.UploadFileHandler)))
// 	http.Handle("/delete", middleware.AuthMiddleware(http.HandlerFunc(handlers.DeleteFileHandler)))

// 	log.Fatal(http.ListenAndServe(":8080", nil))
// }s

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

// var r = mux.NewRouter()

// func main() {

// 	utils.InitRedis()

// 	// Load environment variables from .env file
// 	err := godotenv.Load()
// 	if err != nil {
// 		log.Fatalf("Error loading .env file: %v", err)
// 	}

// 	// Retrieve JWT secret from environment variable
// 	jwtSecret = os.Getenv("JWT_SECRET")
// 	if jwtSecret == "" {
// 		log.Fatal("JWT_SECRET environment variable not set")
// 	}

// 	connStr := "user=" + os.Getenv("POSTGRES_USER") + " password=" + os.Getenv("POSTGRES_PASSWORD") + " dbname=file_sharing sslmode=disable"
// 	db, err = sql.Open("postgres", connStr)
// 	if err != nil {
// 		log.Fatalf("Error opening database: %v", err)
// 	}

// 	// Test the connection
// 	err = db.Ping()
// 	if err != nil {
// 		log.Fatalf("Error pinging database: %v", err)
// 	}
// 	log.Println("Database connection established")

// 	// Route configuration
// 	http.Handle("/register", http.HandlerFunc(handlers.RegisterHandler))
// 	http.Handle("/login", http.HandlerFunc(handlers.LoginHandler))
// 	http.Handle("/upload", middleware.AuthMiddleware(http.HandlerFunc(handlers.UploadFileHandler)))
// 	http.Handle("/delete", middleware.AuthMiddleware(http.HandlerFunc(handlers.DeleteFileHandler)))
// 	// http.Handle("/files", middleware.AuthMiddleware(http.HandlerFunc(handlers.GetFilesHandler)))
// 	http.Handle("/share", middleware.AuthMiddleware(http.HandlerFunc(handlers.ShareFileHandler)))
// 	http.Handle("/files", middleware.AuthMiddleware(http.HandlerFunc(handlers.SearchFilesHandler))) // Add this line
// 	log.Fatal(http.ListenAndServe(":8080", nil))
// }

// Functions for search and caching
// func main() {
// 	utils.InitRedis()

// 	// Load environment variables from .env file
// 	err := godotenv.Load()
// 	if err != nil {
// 		log.Fatalf("Error loading .env file: %v", err)
// 	}

// 	// Retrieve JWT secret from environment variable
// 	jwtSecret = os.Getenv("JWT_SECRET")
// 	if jwtSecret == "" {
// 		log.Fatal("JWT_SECRET environment variable not set")
// 	}

// 	connStr := "user=" + os.Getenv("POSTGRES_USER") + " password=" + os.Getenv("POSTGRES_PASSWORD") + " dbname=file_sharing sslmode=disable"
// 	db, err = sql.Open("postgres", connStr)
// 	if err != nil {
// 		log.Fatalf("Error opening database: %v", err)
// 	}

// 	// Test the connection
// 	err = db.Ping()
// 	if err != nil {
// 		log.Fatalf("Error pinging database: %v", err)
// 	}
// 	log.Println("Database connection established")

// 	// Initialize the background job
// 	// jobs
// 	jobs.Init(connStr)
// 	go jobs.RunBackgroundJob(24 * time.Hour) // Run the job every 24 hours

// 	// Router configuration
// 	r := mux.NewRouter()
// 	r.HandleFunc("/register", handlers.RegisterHandler).Methods("POST")
// 	r.HandleFunc("/login", handlers.LoginHandler).Methods("POST")
// 	r.Handle("/upload", middleware.AuthMiddleware(http.HandlerFunc(handlers.UploadFileHandler)))
// 	r.Handle("/delete", middleware.AuthMiddleware(http.HandlerFunc(handlers.DeleteFileHandler)))
// 	r.Handle("/files/{id:[0-9]+}", middleware.AuthMiddleware(http.HandlerFunc(handlers.UpdateFileHandler))).Methods("PUT")
// 	r.Handle("/files", middleware.AuthMiddleware(http.HandlerFunc(handlers.SearchFilesHandler))).Methods("GET")
// 	r.Handle("/share", middleware.AuthMiddleware(http.HandlerFunc(handlers.ShareFileHandler)))

// 	log.Fatal(http.ListenAndServe(":8080", r))
// }

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
