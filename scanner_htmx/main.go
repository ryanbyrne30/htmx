package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ryanbyrne30/htmx/scanner_htmx/router"
)

type config struct {
	publicDir string
}

func main() {
	var cfg config

	flag.StringVar(&cfg.publicDir, "public", "./static", "the directory to server static files from. Defaults to ./static")
	flag.Parse()

	r := mux.NewRouter()
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(cfg.publicDir))))

	router.Init(r, InfoLog, ErrorLog)
	log.Fatal(http.ListenAndServe(":8080", r))
}