package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/rs/cors" // Import the cors package

	"github.com/golang-jwt/jwt"
	_ "github.com/lib/pq"
)

var db *sql.DB
var jwtKey = []byte(os.Getenv("JWT_SECRET"))

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {

	host := os.Getenv("DB_HOST") // Should be "db"
	port := os.Getenv("DB_PORT") // Default: 5432
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error
	db, err = sql.Open("postgres", psqlInfo) // Assign to the global `db` variable
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	createTableSQL := `
    CREATE TABLE IF NOT EXISTS admin_user (
        id SERIAL PRIMARY KEY,
        username VARCHAR(50) NOT NULL UNIQUE,
        email VARCHAR(100) NOT NULL UNIQUE,
        password VARCHAR(100) NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );`

	_, err = db.Exec(createTableSQL)

	if err != nil {
		log.Fatal(err)
	}
	// Read the SQL file using os.ReadFile (since ioutil is deprecated)
	filePath := "schema.sql"
	sqlBytes, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal("Error reading SQL file:", err)
	}

	// Execute the SQL from the file
	_, err = db.Exec(string(sqlBytes))
	if err != nil {
		log.Fatal("Error executing SQL file:", err)
	}

	fmt.Println("Table 'banners' created successfully from SQL file!")

	// Test the database connection
	err = db.Ping()
	if err != nil {
		log.Fatal("Could not connect to the database: ", err)
	}
	fmt.Println("Successfully connected to the database!")

	if !AdminExists(db) {
		fmt.Println("No admin user found, creating admin user...")

		username := os.Getenv("USERNAME")
		email := "admin@example.com"
		password := os.Getenv("PASSWORD")

		err = CreateAdminUser(db, username, email, password)
		if err != nil {
			log.Fatal(err)
		} else {
			fmt.Println("Admin user created successfully")
		}
	} else {
		fmt.Println("Admin user already exists")
	}
	// Initialize your handlers as usual
	mux := http.NewServeMux()
	mux.HandleFunc("/login", loginHandler)
	mux.HandleFunc("/createbanner", authenticate(createBanner))
	mux.HandleFunc("/banners", getBanners)
	mux.HandleFunc("/product/create", authenticate(createProduct))
	mux.HandleFunc("/product/delete", authenticate(deleteProduct))
	// createJsonTest
	mux.HandleFunc("/product/update", authenticate(updateProduct))
	// mux.HandleFunc("/test/json", createJsonTest)
	mux.HandleFunc("/tags/create", authenticate(createTag))
	mux.HandleFunc("/tags", getTags)
	mux.HandleFunc("/tags/update", authenticate(updateTag))
	mux.HandleFunc("/tags/delete", authenticate(deleteTag))
	mux.HandleFunc("/category/create", authenticate(createCategory))
	mux.HandleFunc("/categories", getCategories)
	mux.HandleFunc("/category/update", authenticate(updateCategory))
	mux.HandleFunc("/category/delete", authenticate(deleteCategory))

	// Wrap the mux with CORS handling
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"}, // Replace with your frontend domain
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	// Start the server with the CORS-enabled handler
	fmt.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", c.Handler(mux)))

}

func authenticate(next http.HandlerFunc) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		tokenStr := cookie.Value

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)

	})
}
func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "http://localhost:3000") // Specify your frontend URL
	(*w).Header().Set("Access-Control-Allow-Credentials", "true")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS") // Specify allowed methods

	// (*w).Header().Set("Access-Control-Allow-Origin", "*") // Allow all origins for development; restrict in production
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

// func enableCors(w *http.ResponseWriter) {
// 	(*w).Header().Set("Access-Control-Allow-Origin", "http://localhost:3000") // Specify your frontend URL
// 	(*w).Header().Set("Access-Control-Allow-Credentials", "true")
// 	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS") // Specify allowed methods
// 	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
// }
