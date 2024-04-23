package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	"amencia.net/quotebox/pkg/models/postgresql"
	_ "github.com/lib/pq" // third party package
)

func setUpDB(dsn string) (*sql.DB, error) {

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
	quotes   *postgresql.QuoteModel
	errorLog *log.Logger
	infoLog  *log.Logger
}

// dsn : data source name
func main() {
	//Create a command line flag
	addr := flag.String("addr", ":4000", "HTTP Network Address")
	dsn := flag.String("dsn",
		os.Getenv("QUOTEBOX_DB_DSN"),
		"PosrgreSQL DSN (Data Source Name)")
	flag.Parse()
	//create a logger
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	var db, err = setUpDB(*dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close() // Always do this before exiting
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		quotes: &postgresql.QuoteModel{
			DB: db,
		},
	}

	srv := &http.Server{
		Addr:     *addr,
		Handler:  app.routes(),
		ErrorLog: errorLog,
	}

	infoLog.Printf("Starting server on port %s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}
