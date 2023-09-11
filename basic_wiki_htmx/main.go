package main

import (
	"html/template"
	"log"
	"net/http"
	"regexp"
)

var templates = template.Must(template.ParseFiles("templates/header.html", "templates/edit.html", "templates/view.html", "templates/click.html"))
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

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

func buttonClickHandler(w http.ResponseWriter, r *http.Request) {
	count, err := ReadCount()
	if err != nil {
		handleError(w, err)
	}
	err = templates.ExecuteTemplate(w, "click.html", count)
	if err != nil {
		handleError(w, err)
	}
	err = WriteCount(count + 1)
	if err != nil {
		handleError(w, err)
	}
}

func renderTemplate(w http.ResponseWriter, temp string, page *Page) {
	name := temp + ".html"
	err := templates.ExecuteTemplate(w, name, page)
	if err != nil {
		handleError(w, err)
	}
}

func handleError(w http.ResponseWriter, err error) {
	log.Fatal(err)
	http.Error(w, err.Error(), http.StatusInternalServerError)
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


func main() {
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	http.HandleFunc("/click/", buttonClickHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}