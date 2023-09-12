package router

import (
	"html/template"

	"github.com/gorilla/mux"
	"github.com/ryanbyrne30/htmx/scanner_htmx/context"
)

var countTemplate = template.Must(template.ParseFiles("templates/count.html"))
var count = 0;

func InitApi(r *mux.Router) {
	r.HandleFunc("/api/count/", MakeHandler(countHandler))
}

func countHandler(ctx context.Context) {
	count += 1
	RenderTemplate(ctx.GetResponseWriter(), countTemplate, "count", count)
}