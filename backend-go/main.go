package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var db *sql.DB

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	var errDb error
	db, errDb = sql.Open("postgres", connStr)
	if errDb != nil {
		log.Fatal(errDb)
	}
	defer db.Close()

	errDb = db.Ping()
	if errDb != nil {
		log.Fatal(errDb)
	}

	fmt.Println("Successfully connected to the database!")

	// Create tables if they don't exist
	createTables()

	// Start RSS fetcher
	StartRSSFetcher()

	// Setup router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://localhost:3000", "http://localhost:8000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Routes
	r.Post("/register", Register)
	r.Post("/login", Login)

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(AuthMiddlewareFunc)
		r.Get("/protected", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Welcome, authenticated user!")
		})
		r.Get("/feeds", GetFeeds)
		r.Post("/feeds", AddFeed)
		r.Delete("/feeds/{feedID}", DeleteFeed)
		r.Get("/posts", GetPosts)
	})

	log.Println("Server starting on port 8000...")
	log.Fatal(http.ListenAndServe(":8000", r))
}

func createTables() {
	createUsersTableSQL := `CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL
	);`

	_, err := db.Exec(createUsersTableSQL)
	if err != nil {
		log.Fatalf("Error creating users table: %v", err)
	}

	createFeedsTableSQL := `CREATE TABLE IF NOT EXISTS feeds (
		id SERIAL PRIMARY KEY,
		user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
		name TEXT NOT NULL,
		url TEXT UNIQUE NOT NULL,
		last_fetched TIMESTAMP
	);`

	_, err = db.Exec(createFeedsTableSQL)
	if err != nil {
		log.Fatalf("Error creating feeds table: %v", err)
	}

	createPostsTableSQL := `CREATE TABLE IF NOT EXISTS posts (
		id SERIAL PRIMARY KEY,
		feed_id INTEGER REFERENCES feeds(id) ON DELETE CASCADE,
		title TEXT NOT NULL,
		url TEXT UNIQUE NOT NULL,
		published_at TIMESTAMP
	);`

	_, err = db.Exec(createPostsTableSQL)
	if err != nil {
		log.Fatalf("Error creating posts table: %v", err)
	}

	fmt.Println("Tables created successfully!")
}
