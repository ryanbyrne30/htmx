package main

import (
	"log"
	"net/http"

	"github.com/ryanbyrne30/htmx/basic_wiki_htmx/routes/pages"
)









func main() {
	http.HandleFunc("/view/", pages.ViewHandler)
	http.HandleFunc("/edit/", EditHandler)
	http.HandleFunc("/save/", SaveHandler)
	// http.HandleFunc("/click/", buttonClickHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}