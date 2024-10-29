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

		username := "admin"
		email := "admin@example.com"
		password := "securepassword"

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
	mux.HandleFunc("/createbanner", createBanner)
	mux.HandleFunc("/banners", getBanners)
	mux.HandleFunc("/createproduct", createProduct)
	mux.HandleFunc("/deleteproduct", deleteProduct)
	mux.HandleFunc("/updateproduct", updateProduct)
	mux.HandleFunc("/tags/create", createTag)
	mux.HandleFunc("/tags", getTags)
	mux.HandleFunc("/tags/update", updateTag)
	mux.HandleFunc("/tags/delete", deleteTag)
	mux.HandleFunc("/category/create", createCategory)
	mux.HandleFunc("/categories", getCategories)
	mux.HandleFunc("/category/update", updateCategory)
	mux.HandleFunc("/category/delete", deleteCategory)

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

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*") // Allow all origins for development; restrict in production
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
}
