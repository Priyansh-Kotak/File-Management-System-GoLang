package jobs

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"time"

	_ "github.com/lib/pq"
)

const (
    checkInterval   = 10 * time.Second // Check interval for expired files
    expirationTime  = 1 * time.Minute  // Expiration time for files
)

var db *sql.DB

// Init initializes the database connection for the background job
func Init(connStr string) {
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
}

// RunBackgroundJob starts the background worker for file expiration
func RunBackgroundJob(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			deleteExpiredFiles()
		}
	}
}

func deleteExpiredFiles() {
	now := time.Now()

	// Get expired files from the database
	rows, err := db.Query("SELECT id, file_name FROM files WHERE upload_date < $1", now.Add(-expirationTime).Format("2006-01-02 15:04:05"))
	if err != nil {
		log.Printf("Error querying expired files: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var fileName string
		if err := rows.Scan(&id, &fileName); err != nil {
			log.Printf("Error scanning expired file: %v", err)
			continue	
		}

		// Delete file from local storage
		filePath := filepath.Join("upload", fileName)
		if err := os.Remove(filePath); err != nil {
			log.Printf("Error deleting file %s: %v", filePath, err)
		} else {
			log.Printf("Deleted file %s", filePath)
		}

		// Delete file metadata from the database
		_, err := db.Exec("DELETE FROM files WHERE id = $1", id)
		if err != nil {
			log.Printf("Error deleting metadata for file ID %d: %v", id, err)
		} else {
			log.Printf("Deleted metadata for file ID %d", id)
		}
	}
}
