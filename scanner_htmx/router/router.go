package router

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func Init(r *mux.Router, infoLogger *log.Logger, errorLogger *log.Logger) {
	r.Use(func (next http.Handler) http.Handler {
		return logMiddleware(next, infoLogger)
	})
	InitPages(r)
	InitApi(r)
}

func logMiddleware(next http.Handler, infoLogger *log.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		infoLogger.Printf("%s: %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}