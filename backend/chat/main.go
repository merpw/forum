package main

import (
	"backend/chat/handlers"
	"backend/chat/ws"
	"database/sql"
	"flag"
	"log"
	"net"
	"net/http"
	"time"
)

// TODO: make reusable. Improve logging and add SSL support

// TODO: improve error handling, add response 'error' when needed

func main() {
	log.SetFlags(log.Lshortfile)
	port := flag.String("port", "8081", "specify server port")
	dbFile := flag.String("db", "chat.db", "specify custom database file path")

	flag.Parse()

	listen, err := net.Listen("tcp", ":"+*port)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Server started on http://localhost:%v\n", *port)

	// TODO: add database protection
	db, err := sql.Open("sqlite3", *dbFile+"?_foreign_keys=true") // enable foreign keys
	if err != nil {
		log.Fatal(err)
	}

	h := handlers.New(db)
	err = h.DB.InitDatabase()
	if err != nil {
		log.Fatal(err)
	}

	hub := ws.NewHub(h.PrimaryHandler())
	h.Hub = hub

	http.HandleFunc("/ws", hub.UpgradeHandler)

	httpServer := http.Server{
		Handler:           http.DefaultServeMux,
		ReadHeaderTimeout: 3 * time.Second,
	}

	err = httpServer.Serve(listen)
	if err != nil {
		log.Fatal(err)
	}
}
