package server

import (
	"encoding/json"
	"github.com/jtejido/ccvapi/validation"
	"net/http"
)

const (
	ApiBasePath      = "/card/api/"
	NumberVerifyPath = ApiBasePath + "verify"
)

// Request struct
type Request struct {
	PAN string
}

// Response struct
type Response struct {
	Valid        bool
	Issuer       string
	Error        string
	PatternMatch int
	LengthMatch  int
}

func apiHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			postHandler(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	rq := Request{}

	err := json.NewDecoder(r.Body).Decode(&rq)
	if err != nil {
		panic(err)
	}

	result := validation.Validate(rq.PAN) // get the type first and see if it passes
	w.Header().Set("Content-Type", "application/json")
	res, err := json.Marshal(&Response{Valid: true, Issuer: result.Name, PatternMatch: result.PatternMatch, LengthMatch: result.LengthMatch, Error: result.Error})
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
