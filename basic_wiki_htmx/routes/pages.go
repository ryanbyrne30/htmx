package main

import (
	"net/http"
	"os"
	"regexp"
)

var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

type Page struct {
	Title string
	Body []byte
}

var ViewHandler = makeHandler(viewHandler)
var EditHandler = makeHandler(editHandler)
var SaveHandler = makeHandler(saveHandler)

func (p *Page) save() error {
	filename := "pages/" + p.Title + ".txt"
	return os.WriteFile(filename, p.Body, 0600)
}


func loadPage(title string) (*Page, error) {
	filename := "pages/" + title + ".txt"
	body, err := os.ReadFile(filename)

	if err != nil {
		return nil, err
	}

	return &Page{Title: title, Body: body}, nil
}

func makeHandler(fn func (http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil || len(m) < 3 {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}


func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	page, err := loadPage(title)

	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	} else {
		renderTemplate(w, "view", page)
	}
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	page, err := loadPage(title)
	if err != nil {
		page = &Page{ Title: title }
	} 
	renderTemplate(w, "edit", page)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		handleError(w, err)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}