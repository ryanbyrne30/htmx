package context

import (
	"net/http"
)

type Context interface {
	GetResponseWriter() http.ResponseWriter
	GetRequest() *http.Request
}

type context struct {
	w http.ResponseWriter
	r *http.Request
}

func NewContext(w http.ResponseWriter, r *http.Request) Context {
	return &context { w, r }
}

func (c *context) GetResponseWriter() http.ResponseWriter {
	return c.w
}

func (c *context) GetRequest() *http.Request {
	return c.r
}