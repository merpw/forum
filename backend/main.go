package main

import (
	"flag"
	"forum/server"
	"log"
	"net"
	"net/http"
)

func main() {
	log.SetFlags(log.Lshortfile)

	port := flag.String("port", "8080", "specify server port")
	flag.Parse()

	listen, err := net.Listen("tcp", ":"+*port)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Server started on http://localhost:%v\n", *port)
	err = http.Serve(listen, server.Start())
	if err != nil {
		log.Fatal(err)
	}
}
