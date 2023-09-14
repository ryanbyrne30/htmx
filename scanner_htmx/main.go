package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ryanbyrne30/htmx/scanner_htmx/router"
)

type config struct {
	publicDir string
	port int
}

func main() {
	var cfg config

	flag.StringVar(&cfg.publicDir, "public", "./static", "the directory to server static files from. Defaults to ./static")
	flag.IntVar(&cfg.port, "port", 8080, "port to run the server on. Default 8080")
	flag.Parse()

	r := mux.NewRouter()
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(cfg.publicDir))))

	router.Init(r, infoLog, errorLog)

	srv := &http.Server{
		Addr: fmt.Sprintf(":%v", cfg.port),
		ErrorLog: errorLog,
		Handler: r,
	}

	infoLog.Printf("Starting server on port %v", cfg.port)

	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}