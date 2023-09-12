package router

import (
	"html/template"

	"github.com/gorilla/mux"
	"github.com/ryanbyrne30/htmx/scanner_htmx/context"
)

var homeTemplates = template.Must(template.ParseFiles("templates/home.html", "templates/base.html"))
var counterTemplates = template.Must(template.ParseFiles("templates/counter.html", "templates/base.html"))

func InitPages(r *mux.Router) {
	r.HandleFunc("/count/", MakeHandler(counterHandler))
	r.HandleFunc("/", MakeHandler(homeHandler))
}

func renderPageTemplate(ctx context.Context, templates *template.Template, data any) {
	RenderTemplate(ctx.GetResponseWriter(), templates, "base", data)
}

func homeHandler(ctx context.Context) {
	renderPageTemplate(ctx, homeTemplates, nil)
}

func counterHandler(ctx context.Context) {
	renderPageTemplate(ctx, counterTemplates, 0)
}
