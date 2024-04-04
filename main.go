package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"amencia.net/quotebox/pkg/postgresql"
	_ "github.com/lib/pq" // third party package
)

func setUpDB() (*sql.DB, error) {

	// Provide the credentials for our database
	const (
		host     = "localhost"
		port     = 5432
		user     = "quotebox"
		password = "iamme"
		dbname   = "quotebox"
	)
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	// Establish a connection to the database
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	// Test our connection
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

// dependecies (things/variables)
type application struct {
	quotes *postgresql.QuoteModel
}

// dsn : data source name
func main() {
	var db, err = setUpDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close() // Always do this before exiting
	app := &application{
		quotes: &postgresql.QuoteModel{
			DB: db,
		},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/quote", app.createQuoteForm)
	mux.HandleFunc("/quote-add", app.createQuote)
	//mux.HandleFunc("/show", app.displayQuotation)
	log.Println("Starting server on port :4000")
	err = http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
