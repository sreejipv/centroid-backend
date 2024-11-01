package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	// Import the cors package

	"github.com/golang-jwt/jwt"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

func AdminExists(db *sql.DB) bool {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM admin_user").Scan(&count)
	if err != nil {
		log.Fatal("Error checking admin user existence:", err)
	}
	return count > 0
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CreateAdminUser(db *sql.DB, username, email, password string) error {
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return err
	}

	_, err = db.Exec("INSERT INTO admin_user (username, email, password) VALUES ($1, $2, $3)", username, email, hashedPassword)
	return err
}

// loginHandler handles the admin login
func loginHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var loginRequest LoginRequest
	err := json.NewDecoder(r.Body).Decode(&loginRequest)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	username := loginRequest.Username
	password := loginRequest.Password

	if username == "" || password == "" {
		http.Error(w, "Username and password required", http.StatusBadRequest)
		return
	}

	// Check if the database is connected
	err = db.Ping()
	if err != nil {
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	var dbUsername, dbPasswordHash string
	err = db.QueryRow("SELECT username, password FROM admin_user WHERE username = $1", username).Scan(&dbUsername, &dbPasswordHash)
	if err != nil {
		http.Error(w, "Invalid user", http.StatusUnauthorized)
		return
	}

	// Compare the hashed password
	err = bcrypt.CompareHashAndPassword([]byte(dbPasswordHash), []byte(password))
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  expirationTime,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode, // Adjust as necessary
		Secure:   false,                // Set to true in production with HTTPS
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // Set response status to 200
	json.NewEncoder(w).Encode(map[string]string{
		"message": fmt.Sprintf("Welcome, %s!", dbUsername),
	})
}
