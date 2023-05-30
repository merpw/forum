package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
)

// ErrorResponse responses with specified error code in format "404 Not Found"
func ErrorResponse(w http.ResponseWriter, code int) {
	http.Error(w, fmt.Sprintf("%v %v", code, http.StatusText(code)), code)
}

// SendObject sends object to http.ResponseWriter
//
// panics if error occurs
func SendObject(w http.ResponseWriter, object any) {
	w.Header().Set("Content-Type", "application/json")
	objJson, err := json.Marshal(object)
	if err != nil {
		log.Panic(err)
		return
	}
	_, err = w.Write(objJson)
	if err != nil {
		log.Panic(err)
		return
	}
}

type Route struct {
	Method  string
	Pattern *regexp.Regexp
	Handler http.HandlerFunc
}

func NewRoute(method, pattern string, handler http.HandlerFunc) Route {
	return Route{method, regexp.MustCompile("^" + pattern + "$"), handler}
}
