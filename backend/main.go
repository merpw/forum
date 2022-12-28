package main

import (
	"database/sql"
	"flag"
	"forum/server"
	"log"
	"net"
	"net/http"
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

	db, err := sql.Open("sqlite3", *dbFile+"?_foreign_keys=true") // enable foreign keys
	if err != nil {
		log.Fatal(err)
	}

	srv := server.Connect(db)
	err = srv.DB.InitDatabase()
	if err != nil {
		log.Fatal(err)
	}

	err = http.Serve(listen, srv.Start())
	if err != nil {
		log.Fatal(err)
	}
}
