package server

import (
	"forum/server"
	"net/http/httptest"
	"testing"
)

func TestStart(t *testing.T) {
	router := server.Start()
	srv := httptest.NewServer(router)
	defer srv.Close()
}
