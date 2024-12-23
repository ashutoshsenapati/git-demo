// main.go
package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "password"
	dbname   = "testdb"
)

type Book struct {
	ID        int
	Title     string
	Author    string
	Published time.Time
}

func initDB() (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %v", err)
	}

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

func insertBook(db *sql.DB, book Book) error {
	query := `
        INSERT INTO books (title, author, published)
        VALUES ($1, $2, $3)
        RETURNING id`

	return db.QueryRow(query, book.Title, book.Author, book.Published).Scan(&book.ID)
}

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
	db, err := initDB()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

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
