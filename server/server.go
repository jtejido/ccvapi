package server

import (
	"github.com/jtejido/ccvapi/logging"
	"net/http"
)

type ServeMux struct {
	Logger logging.Logging
	mux    *http.ServeMux
}

func NewServeMux(options ...func(*ServeMux)) *ServeMux {
	s := &ServeMux{
		mux: http.NewServeMux(),
	}

	for _, f := range options {
		f(s)
	}

	s.mux.Handle(NumberVerifyPath, loader(apiHandler(), logIt(s.Logger)))
	return s
}

func (s *ServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}
