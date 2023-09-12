package router

import (
	"html/template"
	"log"
	"net/http"

	"github.com/ryanbyrne30/htmx/scanner_htmx/context"
)

func MakeHandler(handler func (ctx context.Context)) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		ctx := context.NewContext(w, r)
		handler(ctx)
	}
}

func RenderTemplate(w http.ResponseWriter, templates *template.Template, temp string, data any) {
	name := temp + ".html"
	err := templates.ExecuteTemplate(w, name, data)
	if err != nil {
		handleError(w, err)
	}
}

func handleError(w http.ResponseWriter, err error) {
	log.Fatal(err)
	http.Error(w, err.Error(), http.StatusInternalServerError)
}