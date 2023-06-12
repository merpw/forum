package main

import (
	"backend/forum/handlers"
	"database/sql"
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

func main() {
	log.SetFlags(log.Lshortfile)

	if os.Getenv("FORUM_BACKEND_SECRET") == "" {
		log.Println("WARNING: FORUM_BACKEND_SECRET is not set, any request with Internal-Auth header will be accepted")
	}

	if os.Getenv("FORUM_IS_PRIVATE") == "" {
		err := os.Setenv("FORUM_IS_PRIVATE", "true")
		if err != nil {
			log.Fatal(err)
		}
		log.Println("INFO: FORUM_IS_PRIVATE is not set, defaulting to true")
	}

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

	h := handlers.New(db)
	err = h.DB.InitDatabase()
	if err != nil {
		log.Fatal(err)
	}

	httpServer := http.Server{
		Handler:           h.Handler(),
		ReadHeaderTimeout: 3 * time.Second,
	}

	err = httpServer.Serve(listen)
	if err != nil {
		log.Fatal(err)
	}
}
