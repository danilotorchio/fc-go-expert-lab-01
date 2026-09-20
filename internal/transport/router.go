package transport

import (
	"net/http"

	"github.com/danilotorchio/fc-go-expert-lab-01/internal/weather"
)

type Router struct {
	mux *http.ServeMux
}

func (r *Router) Handler(handler *weather.Handler) http.Handler {
	r.mux.Handle("GET /weather/{cep}", handler)
	return r.mux
}

func NewRouter() *Router {
	return &Router{
		mux: http.NewServeMux(),
	}
}
