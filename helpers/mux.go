package helpers

import (
	"net/http"
)

// GeMux -
type GeMux struct {
	http.ServeMux
	middlewares []func(http.Handler) http.Handler
}

func (gm *GeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var h http.Handler = &gm.ServeMux

	for _, next := range gm.middlewares {
		h = next(h)
	}

	h.ServeHTTP(w, r)
}

// Use -
func (gm *GeMux) Use(middleware func(next http.Handler) http.Handler) {
	gm.middlewares = append(gm.middlewares, middleware)
}
