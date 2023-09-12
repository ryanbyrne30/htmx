package router

import (
	"html/template"

	"github.com/gorilla/mux"
	"github.com/ryanbyrne30/htmx/scanner_htmx/context"
)

var templates = template.Must(template.ParseFiles("templates/home.html"))
var homeTemplates = template.Must(template.ParseFiles("templates/home.html", "templates/base.html"))

func InitPages(r *mux.Router) {
	r.HandleFunc("/", MakeHandler(homeHandler))
}

func homeHandler(ctx context.Context) {
	RenderTemplate(ctx.GetResponseWriter(), homeTemplates, "base", nil)
}
