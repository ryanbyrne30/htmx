package router

import "github.com/gorilla/mux"

func Init(r *mux.Router) {
	InitPages(r)
	InitApi(r)
}

