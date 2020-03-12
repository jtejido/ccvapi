package server

import (
	"github.com/jtejido/ccvapi/logging"
	"net/http"
	"time"
)

type service func(http.Handler) http.Handler

func loader(h http.Handler, svcs ...service) http.Handler {
	for _, svc := range svcs {
		h = svc(h)
	}
	return h
}

func logIt(a logging.Logging) service {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			h.ServeHTTP(w, r)
			a.Printf("%s request to %s took %vns.", r.Method, r.URL.Path, time.Now().Sub(start).Nanoseconds())
		})
	}
}
