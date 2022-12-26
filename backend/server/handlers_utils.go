package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// errorResponse responses with specified error code in format "404 Not Found"
func errorResponse(w http.ResponseWriter, code int) {
	http.Error(w, fmt.Sprintf("%v %v", code, http.StatusText(code)), code)
}

func (srv *Server) methodGetHandler(w http.ResponseWriter, r *http.Request) {}

func (srv *Server) methodPostHandler(w http.ResponseWriter, r *http.Request) {}

// sendObject sends object to http.ResponseWriter
//
// calls errorResponse(500) if error happened
func sendObject(w http.ResponseWriter, object any) {
	w.Header().Set("Content-Type", "application/json")
	objJson, err := json.Marshal(object)
	if err != nil {
		log.Println(err)
		errorResponse(w, 500)
		return
	}
	_, err = w.Write(objJson)
	if err != nil {
		log.Println(err)
		errorResponse(w, 500)
		return
	}
}
