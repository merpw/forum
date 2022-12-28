package main

import (
	"flag"
	"forum/database"
	"forum/server"
	"io"
	"log"
	"net"
	"net/http"
	"os"
)

func main() {
	log.SetFlags(log.Lshortfile)

	// define multiplex output for log messages to log.txt and stdout
	logFile, err := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()
	multi := io.MultiWriter(logFile, os.Stdout)
	log.SetOutput(multi)

	port := flag.String("port", "8080", "specify server port")
	flag.Parse()

	listen, err := net.Listen("tcp", ":"+*port)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Server started on http://localhost:%v\n", *port)
	err = database.InitDatabase()
	if err != nil {
		log.Fatal(err)
	}
	err = http.Serve(listen, server.Start())
	if err != nil {
		log.Fatal(err)
	}
}
