package main

import (
	"database/sql"
	"flag"
	"forum/server"
	"log"
	"net"
	"net/http"
	"time"
)

func main() {
	log.SetFlags(log.Lshortfile)

	port := flag.String("port", "8080", "specify server port")
	dbFile := flag.String("db", "database.db", "specify custom database file path")

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

	srv := server.Connect(db)
	err = srv.DB.InitDatabase()
	if err != nil {
		log.Fatal(err)
	}

	httpServer := http.Server{
		Handler:           srv.Start(),
		ReadHeaderTimeout: 3 * time.Second,
	}

	err = httpServer.Serve(listen)
	if err != nil {
		log.Fatal(err)
	}
}
