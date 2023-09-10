package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

type Page struct {
	Title string
	Body []byte
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return os.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := os.ReadFile(filename)

	if err != nil {
		return nil, err
	}

	return &Page{Title: filename, Body: body}, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view/"):]
	page, err := loadPage(title)

	if err != nil {
		fmt.Fprintf(w, "<h1>Not Found</h1>")
	} else {
		fmt.Fprintf(w, "<h1>%s</h1><main>%s</main>", page.Title, page.Body)
	}

}

func main() {
	http.HandleFunc("/view/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}