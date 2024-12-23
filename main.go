// main.go
package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

// Database connection parameters
const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "password"
	dbname   = "testdb"
)

// Book represents a row in our books table
type Book struct {
	ID        int
	Title     string
	Author    string
	Published time.Time
}

// initDB establishes the database connection and creates the books table
func initDB() (*sql.DB, error) {
	// Construct the connection string
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Open a connection to the database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	// Verify the connection is working
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %v", err)
	}

	// Create the books table if it doesn't exist
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS books (
            id SERIAL PRIMARY KEY,
            title VARCHAR(100) NOT NULL,
            author VARCHAR(100) NOT NULL,
            published DATE NOT NULL
        )
    `)
	if err != nil {
		return nil, fmt.Errorf("error creating table: %v", err)
	}

	return db, nil
}

// insertBook adds a new book to the database
func insertBook(db *sql.DB, book Book) error {
	query := `
        INSERT INTO books (title, author, published)
        VALUES ($1, $2, $3)
        RETURNING id`

	return db.QueryRow(query, book.Title, book.Author, book.Published).Scan(&book.ID)
}

// getBookByID retrieves a book from the database by its ID
func getBookByID(db *sql.DB, id int) (Book, error) {
	var book Book
	query := `
        SELECT id, title, author, published
        FROM books
        WHERE id = $1`

	err := db.QueryRow(query, id).Scan(&book.ID, &book.Title, &book.Author, &book.Published)
	if err != nil {
		return Book{}, fmt.Errorf("error getting book: %v", err)
	}

	return book, nil
}

// getAllBooks retrieves all books from the database
func getAllBooks(db *sql.DB) ([]Book, error) {
	query := `
        SELECT id, title, author, published
        FROM books`

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying books: %v", err)
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var book Book
		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Published)
		if err != nil {
			return nil, fmt.Errorf("error scanning book: %v", err)
		}
		books = append(books, book)
	}

	return books, nil
}

func main() {
	// Initialize the database connection
	db, err := initDB()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// Insert a sample book
	newBook := Book{
		Title:     "The Go Programming Language",
		Author:    "Alan A. A. Donovan and Brian W. Kernighan",
		Published: time.Date(2015, 11, 1, 0, 0, 0, 0, time.UTC),
	}

	err = insertBook(db, newBook)
	if err != nil {
		log.Fatal("Failed to insert book:", err)
	}
	fmt.Println("Successfully inserted new book")

	// Retrieve and display all books
	books, err := getAllBooks(db)
	if err != nil {
		log.Fatal("Failed to get books:", err)
	}

	fmt.Println("\nAll books in database:")
	for _, book := range books {
		fmt.Printf("ID: %d, Title: %s, Author: %s, Published: %s\n",
			book.ID, book.Title, book.Author, book.Published.Format("2006-01-02"))
	}
}
