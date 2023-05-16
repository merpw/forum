package server_test

import (
	"database/sql"
	"forum/server"
	"log"
	"net/http/httptest"
	"os"
	"testing"
)

var testServer *httptest.Server

// TestMain is the entry point for all tests
func TestMain(m *testing.M) {
	testServer = startServer()
	os.Exit(m.Run())
}

// startServer starts a test server with new test database file and returns a httptest.Server
func startServer() *httptest.Server {
	// This code is attempting to remove a file named "test.db" from the current
	// directory. If the file does not exist, it will not return an error. If an error
	// occurs during the removal process, it will check if the error is due to the
	// file not existing. If the error is not due to the file not existing, it will
	// log the error and exit the program.
	err := os.Remove("./test.db")
	if err != nil {
		if !os.IsNotExist(err) {
			log.Fatal(err)
		}
	}

	// This code is creating a new SQLite database file named "test.db" with foreign
	// key support enabled. It then opens a connection to the database using the
	// `sql.Open()` function from the `database/sql` package. If an error occurs during
	// the opening of the database connection, it will log the error and exit the
	// program.
	db, err := sql.Open("sqlite3", "./test.db?_foreign_keys=true")
	if err != nil {
		log.Fatal(err)
	}

	// This code is creating a new server instance using the `server.Connect()`
	srv := server.Connect(db)

	// This code is initializing the database using the `srv.DB.InitDatabase()`
	err = srv.DB.InitDatabase()
	if err != nil {
		log.Fatal(err)
	}

	// This code is starting the server using the `srv.Start()` method and wrapping
	// the returned `http.Handler` in a `httptest.NewServer()` instance.
	return httptest.NewServer(srv.Start())
}

// giveMeUsrPlusCookie returns a userid int and cookie string
// func giveMeUsrPlusCookie() (int, string) {
// 	// create a struct for the user
// 	tuser1 := TestUser{
// 		"tuser1",
// 		"tuser1@temail1.com",
// 		"tuser1Password",
// 		"tFN1",
// 		"tLN1",
// 		"2000-01-01",
// 		"male",
// 	}

// 	// initiate the database and start the httptest server
// 	tsrv1 := startServer()

// 	// signup the user and get the response int (userid)

// 	// signin the user and get the response string (cookie)

// 	// return the userid and cookie

// }
