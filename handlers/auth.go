// // // package handlers

// // // import (
// // // 	"encoding/json"
// // // 	"file-management/models"
// // // 	"file-management/utils"
// // // 	"net/http"

// // // 	"golang.org/x/crypto/bcrypt"
// // // )

// // // // Mock user database
// // // var userDB = make(map[string]models.User)

// // // // RegisterHandler handles user registration and JWT token generation
// // // func RegisterHandler(w http.ResponseWriter, r *http.Request) {
// // // 	if r.Method != http.MethodPost {
// // // 		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
// // // 		return
// // // 	}

// // // 	var user models.User
// // // 	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
// // // 		http.Error(w, "Invalid request payload", http.StatusBadRequest)
// // // 		return
// // // 	}

// // // 	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
// // // 	if err != nil {
// // // 		http.Error(w, "Error hashing password", http.StatusInternalServerError)
// // // 		return
// // // 	}
// // // 	user.Password = string(hashedPassword)

// // // 	// Store user in mock database
// // // 	userDB[user.Email] = user

// // // 	// Generate JWT token
// // // 	token, err := utils.GenerateJWT(user)
// // // 	if err != nil {
// // // 		http.Error(w, "Error generating token", http.StatusInternalServerError)
// // // 		return
// // // 	}

// // // 	w.Header().Set("Content-Type", "application/json")
// // // 	json.NewEncoder(w).Encode(map[string]string{"token": token})
// // // }

// // // // LoginHandler handles user login and JWT token generation
// // // func LoginHandler(w http.ResponseWriter, r *http.Request) {
// // // 	if r.Method != http.MethodPost {
// // // 		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
// // // 		return
// // // 	}

// // // 	var user models.User
// // // 	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
// // // 		http.Error(w, "Invalid request payload", http.StatusBadRequest)
// // // 		return
// // // 	}

// // // 	storedUser, exists := userDB[user.Email]
// // // 	if !exists || bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password)) != nil {
// // // 		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
// // // 		return
// // // 	}

// // // 	token, err := utils.GenerateJWT(storedUser)
// // // 	if err != nil {
// // // 		http.Error(w, "Error generating token", http.StatusInternalServerError)
// // // 		return
// // // 	}

// // // 	w.Header().Set("Content-Type", "application/json")
// // // 	json.NewEncoder(w).Encode(map[string]string{"token": token})
// // // }

// // // // ProtectedHandler demonstrates token authentication
// // // func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
// // // 	tokenString := r.Header.Get("Authorization")
// // // 	claims, err := utils.VerifyJWT(tokenString)
// // // 	if err != nil {
// // // 		http.Error(w, "Unauthorized", http.StatusUnauthorized)
// // // 		return
// // // 	}

// // // 	userEmail, ok := claims["email"].(string)
// // // 	if !ok {
// // // 		http.Error(w, "Unauthorized", http.StatusUnauthorized)
// // // 		return
// // // 	}

// // // 	user, exists := userDB[userEmail]
// // // 	if !exists {
// // // 		http.Error(w, "User not found", http.StatusUnauthorized)
// // // 		return
// // // 	}

// // // 	w.Header().Set("Content-Type", "application/json")
// // // 	json.NewEncoder(w).Encode(user)
// // // }

// // package handlers

// // import (
// // 	"database/sql"
// // 	"encoding/json"
// // 	"file-management/utils"
// // 	"net/http"

// // 	_ "github.com/lib/pq"
// // 	"golang.org/x/crypto/bcrypt"
// // )

// // var db *sql.DB // Initialize your database connection here

// // type RegisterRequest struct {
// // 	Email    string `json:"email"`
// // 	Password string `json:"password"`
// // }

// // type LoginRequest struct {
// // 	Email    string `json:"email"`
// // 	Password string `json:"password"`
// // }

// // // RegisterHandler handles user registration
// // func RegisterHandler(w http.ResponseWriter, r *http.Request) {
// // 	if r.Method != http.MethodPost {
// // 		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
// // 		return
// // 	}

// // 	var req RegisterRequest
// // 	err := json.NewDecoder(r.Body).Decode(&req)
// // 	if err != nil {
// // 		http.Error(w, "Invalid request body", http.StatusBadRequest)
// // 		return
// // 	}

// // 	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
// // 	if err != nil {
// // 		http.Error(w, "Error hashing password", http.StatusInternalServerError)
// // 		return
// // 	}

// // 	_, err = db.Exec("INSERT INTO users (email, password_hash) VALUES ($1, $2)", req.Email, hashedPassword)
// // 	if err != nil {
// // 		http.Error(w, "Error registering user", http.StatusInternalServerError)
// // 		return
// // 	}

// // 	w.WriteHeader(http.StatusOK)
// // 	w.Write([]byte("User registered successfully"))
// // }

// // // LoginHandler handles user login
// // func LoginHandler(w http.ResponseWriter, r *http.Request) {
// // 	if r.Method != http.MethodPost {
// // 		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
// // 		return
// // 	}

// // 	var req LoginRequest
// // 	err := json.NewDecoder(r.Body).Decode(&req)
// // 	if err != nil {
// // 		http.Error(w, "Invalid request body", http.StatusBadRequest)
// // 		return
// // 	}

// // 	var storedPasswordHash string
// // 	err = db.QueryRow("SELECT password_hash FROM users WHERE email = $1", req.Email).Scan(&storedPasswordHash)
// // 	if err != nil {
// // 		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
// // 		return
// // 	}

// // 	err = bcrypt.CompareHashAndPassword([]byte(storedPasswordHash), []byte(req.Password))
// // 	if err != nil {
// // 		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
// // 		return
// // 	}

// // 	token, err := utils.GenerateJWT(req.Email)
// // 	if err != nil {
// // 		http.Error(w, "Error generating token", http.StatusInternalServerError)
// // 		return
// // 	}

// // 	response := map[string]string{"token": token}
// // 	w.Header().Set("Content-Type", "application/json")
// // 	json.NewEncoder(w).Encode(response)
// // }

// // package handlers

// // import (
// // 	"database/sql"
// // 	"encoding/json"
// // 	"log"
// // 	"net/http"
// // 	"os"
// // 	"time"

// // 	"github.com/dgrijalva/jwt-go"
// // 	"golang.org/x/crypto/bcrypt"
// // )

// // var db *sql.DB
// // var jwtSecret string

// // func init() {
// // 	// var err error
// // 	jwtSecret = os.Getenv("JWT_SECRET")
// // 	if jwtSecret == "" {
// // 		log.Fatal("JWT_SECRET environment variable not set")
// // 	}
// // }

// // func RegisterHandler(w http.ResponseWriter, r *http.Request) {
// // 	if r.Method != http.MethodPost {
// // 		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
// // 		return
// // 	}

// // 	var user struct {
// // 		Email    string `json:"email"`
// // 		Password string `json:"password"`
// // 	}

// // 	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
// // 		http.Error(w, "Invalid request payload", http.StatusBadRequest)
// // 		return
// // 	}

// // 	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
// // 	if err != nil {
// // 		http.Error(w, "Error hashing password", http.StatusInternalServerError)
// // 		return
// // 	}

// // 	_, err = db.Exec("INSERT INTO users (email, password_hash) VALUES ($1, $2)", user.Email, hash)
// // 	if err != nil {
// // 		http.Error(w, "Error registering user", http.StatusInternalServerError)
// // 		return
// // 	}

// // 	w.WriteHeader(http.StatusOK)
// // }

// // func LoginHandler(w http.ResponseWriter, r *http.Request) {
// // 	if r.Method != http.MethodPost {
// // 		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
// // 		return
// // 	}

// // 	var user struct {
// // 		Email    string `json:"email"`
// // 		Password string `json:"password"`
// // 	}

// // 	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
// // 		http.Error(w, "Invalid request payload", http.StatusBadRequest)
// // 		return
// // 	}

// // 	var storedHash string
// // 	err := db.QueryRow("SELECT password_hash FROM users WHERE email = $1", user.Email).Scan(&storedHash)
// // 	if err != nil {
// // 		http.Error(w, "Error retrieving user", http.StatusInternalServerError)
// // 		return
// // 	}

// // 	err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(user.Password))
// // 	if err != nil {
// // 		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
// // 		return
// // 	}

// // 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
// // 		"email": user.Email,
// // 		"exp":   time.Now().Add(time.Hour * 24).Unix(),
// // 	})

// // 	tokenString, err := token.SignedString([]byte(jwtSecret))
// // 	if err != nil {
// // 		http.Error(w, "Error generating token", http.StatusInternalServerError)
// // 		return
// // 	}

// // 	response := map[string]string{"token": tokenString}
// // 	w.Header().Set("Content-Type", "application/json")
// // 	json.NewEncoder(w).Encode(response)
// // }

// package handlers

// import (
// 	"database/sql"
// 	"encoding/json"
// 	"log"
// 	"net/http"
// 	"os"
// 	"time"
// 	"file-management/utils"
// 	"io/ioutil"

// 	"github.com/dgrijalva/jwt-go"
// 	"golang.org/x/crypto/bcrypt"
// )

// var db *sql.DB

// func init() {
// 	var err error
// 	connStr := "user=postgres password=yourpassword dbname=file_sharing sslmode=disable"
// 	db, err = sql.Open("postgres", connStr)
// 	if err != nil {
// 		panic("Error connecting to the database: " + err.Error())
// 	}
// }

// type RegisterRequest struct {
// 	Email    string `json:"email"`
// 	Password string `json:"password"`
// }

// func RegisterHandler(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodPost {
// 		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
// 		return
// 	}

// 	var req RegisterRequest
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		http.Error(w, "Invalid request payload", http.StatusBadRequest)
// 		return
// 	}

// 	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
// 	if err != nil {
// 		http.Error(w, "Error hashing password", http.StatusInternalServerError)
// 		return
// 	}

// 	_, err = db.Exec("INSERT INTO users (email, password_hash) VALUES ($1, $2)", req.Email, hash)
// 	if err != nil {
// 		http.Error(w, "Error saving user", http.StatusInternalServerError)
// 		return
// 	}

// 	// Respond with success
// 	w.WriteHeader(http.StatusOK)
// 	w.Write([]byte("User registered successfully"))
// }

// func LoginHandler(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodPost {
// 		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
// 		return
// 	}

// 	var user struct {
// 		Email    string `json:"email"`
// 		Password string `json:"password"`
// 	}

// 	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
// 		http.Error(w, "Invalid request payload", http.StatusBadRequest)
// 		return
// 	}

// 	var storedHash string
// 	err := db.QueryRow("SELECT password_hash FROM users WHERE email = $1", user.Email).Scan(&storedHash)
// 	if err != nil {
// 		http.Error(w, "Error retrieving user", http.StatusInternalServerError)
// 		return
// 	}

// 	err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(user.Password))
// 	if err != nil {
// 		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
// 		return
// 	}

// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
// 		"email": user.Email,
// 		"exp":   time.Now().Add(time.Hour * 24).Unix(),
// 	})

// 	tokenString, err := token.SignedString([]byte(jwtSecret))
// 	if err != nil {
// 		http.Error(w, "Error generating token", http.StatusInternalServerError)
// 		return
// 	}

// 	response := map[string]string{"token": tokenString}
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(response)
// }

package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"time"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"

	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/dgrijalva/jwt-go"
)

var db *sql.DB
var jwtSecret string

func init() {
	err := godotenv.Load()
	// var err error
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	jwtSecret = os.Getenv("JWT_SECRET")
	log.Println((jwtSecret))
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable not sett")
	}
	connStr := "user=postgres password=priyanshkotak dbname=file_sharing sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		panic("Error connecting to the database: " + err.Error())
	}
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	_, err = db.Exec("INSERT INTO users (email, password_hash) VALUES ($1, $2)", req.Email, hash)
	if err != nil {
		// Log the error and provide more details in the response
		http.Error(w, "Error saving user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with success
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User registered successfully"))
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var user struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var storedHash string
	err := db.QueryRow("SELECT password_hash FROM users WHERE email = $1", user.Email).Scan(&storedHash)
	if err != nil {
		http.Error(w, "Error retrieving user", http.StatusInternalServerError)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(user.Password))
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})

	// jwtSecret := ""
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	response := map[string]string{"token": tokenString}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
