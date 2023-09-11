package main

import (
	"html/template"
	"log"
	"net/http"
)

var templates = template.Must(template.ParseFiles("../templates/header.html", "../templates/edit.html", "../templates/view.html", "../templates/click.html"))


func handleError(w http.ResponseWriter, err error) {
	log.Fatal(err)
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func renderTemplate(w http.ResponseWriter, temp string, page *Page) {
	name := temp + ".html"
	err := templates.ExecuteTemplate(w, name, page)
	if err != nil {
		handleError(w, err)
	}
}