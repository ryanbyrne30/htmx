package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/ryanbyrne30/htmx/scanner_htmx/internal/models"
)

type application struct {
	errorLog *log.Logger
	infoLog *log.Logger
	snippets *models.SnippetModel
}

func main() {
	dsn := flag.String("dsn", "./db/app.db", "port to run the server on. Default 8080")
	port := flag.Int("port", 8080, "port to run the server on. Default 8080")
	flag.Parse()

	// initialize loggers
	var infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	var errorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime)

	// initialize database
	db, err := openDb(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	// initialize application
	app := &application{
		errorLog: errorLog,
		infoLog: infoLog,
		snippets: &models.SnippetModel{DB: db},
	}

	// initialize server
	srv := &http.Server{
		Addr: fmt.Sprintf(":%v", *port),
		ErrorLog: errorLog,
		Handler: app.routes(),
	}

	infoLog.Printf("Starting server on port %v", *port)

	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

func openDb(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
