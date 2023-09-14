package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type config struct {
	publicDir string
	port int
}

type application struct {
	errorLog *log.Logger
	infoLog *log.Logger
}

func main() {
	var cfg config

	flag.StringVar(&cfg.publicDir, "public", "./static", "the directory to server static files from. Defaults to ./static")
	flag.IntVar(&cfg.port, "port", 8080, "port to run the server on. Default 8080")
	flag.Parse()

	var infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	var errorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime)

	app := &application{
		errorLog: errorLog,
		infoLog: infoLog,
	}

	r := mux.NewRouter()
	r.Use(app.logMw)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(cfg.publicDir))))
	r.HandleFunc("/api/count/", app.countClickHandler)
	r.HandleFunc("/count/", app.counterHandler)
	r.HandleFunc("/", app.homeHandler)

	srv := &http.Server{
		Addr: fmt.Sprintf(":%v", cfg.port),
		ErrorLog: errorLog,
		Handler: r,
	}

	infoLog.Printf("Starting server on port %v", cfg.port)

	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}
