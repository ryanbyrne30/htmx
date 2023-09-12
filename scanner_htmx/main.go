package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ryanbyrne30/htmx/scanner_htmx/router"
)

func main() {
	var publicDir string

	flag.StringVar(&publicDir, "public", "./static", "the directory to server static files from. Defaults to ./static")
	flag.Parse()

	r := mux.NewRouter()
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(publicDir))))

	router.Init(r)
	log.Fatal(http.ListenAndServe(":8080", r))
}