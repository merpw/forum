package main

import (
	"backend/common/server"
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"runtime/debug"
	"time"
)

var attachmentsDir *string

func main() {
	log.SetFlags(log.Lshortfile)

	port := flag.String("port", "8082", "specify server port")
	attachmentsDir = flag.String("dir", "./attachments", "specify custom directory to store attachments")

	flag.Parse()

	listen, err := net.Listen("tcp", ":"+*port)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Server started on http://localhost:%v\n", *port)

	fs := http.FileServer(http.Dir(*attachmentsDir))

	err = os.MkdirAll(*attachmentsDir, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/api/attachments/upload", uploadHandler)

	http.Handle("/api/attachments/", http.StripPrefix("/api/attachments/", fs))

	// implement handler with error recovery

	httpServer := http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					log.Printf("ERROR %d. %v\n%s", http.StatusInternalServerError, err, debug.Stack())
					server.ErrorResponse(w, http.StatusInternalServerError) // 500 ERROR
				}
			}()

			http.DefaultServeMux.ServeHTTP(w, r)
		}),
		ReadHeaderTimeout: 3 * time.Second,
	}

	err = httpServer.Serve(listen)
	if err != nil {
		log.Fatal(err)
	}
}
