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

// errorCheck handle basic request errors
func errorBasicCheck(w http.ResponseWriter, r *http.Request, pathToCheck string) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("ERROR %d. %v\n", http.StatusInternalServerError, err)
			errorResponse(w, http.StatusInternalServerError) // 500 ERROR
		}
	}()

	if r.URL.Path != pathToCheck {
		log.Printf("ERROR %d. r.URL.Path = %s != \""+pathToCheck+"\"\n", http.StatusNotFound, r.URL.Path)
		errorResponse(w, http.StatusNotFound) // 404 ERROR
		return
	}

	if err := r.ParseForm(); err != nil {
		log.Printf("ERROR %d. ParseForm() err: %v\n", http.StatusBadRequest, err)
		errorResponse(w, http.StatusBadRequest) // 400 ERROR
		return
	}

	if r.Method != http.MethodPost { // not POST method case
		log.Printf("ERROR %d. %v\n", http.StatusMethodNotAllowed, fmt.Sprintf("request method %s is inappropriate for the URL %s", r.Method, r.URL.Path))
		errorResponse(w, http.StatusMethodNotAllowed)
		return
	}
}

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
