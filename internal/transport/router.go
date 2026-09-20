package transport

import (
	"net/http"
)

type Router struct {
	mux *http.ServeMux
}

func (r *Router) Handler(handler http.Handler) http.Handler {
	r.mux.Handle("GET /weather/{cep}", handler)
	return r.mux
}

func NewRouter() *Router {
	return &Router{
		mux: http.NewServeMux(),
	}
}
