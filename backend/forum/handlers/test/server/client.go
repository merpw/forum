package server

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

// TestClient is like http.Client, but with cookies and some helper methods to make testing easier.
type TestClient struct {
	*http.Client
	testServer *TestServer
	Cookies    []*http.Cookie
	TestClientData
}

// TestClient creates a new TestClient
func (testServer *TestServer) TestClient() *TestClient {
	return &TestClient{
		Client:     testServer.Client(),
		testServer: testServer,
		Cookies:    make([]*http.Cookie, 0),
	}
}

// TestRequest is like http.Client.Do, but with cookies and better testing error messages.
func (cli *TestClient) TestRequest(t testing.TB,
	req *http.Request, responseCode int) (resp *http.Response, respBody []byte) {

	t.Helper()
	for _, cookie := range cli.Cookies {
		req.AddCookie(cookie)
	}

	resp, err := cli.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	respBody, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Response body read error: %v", err)
	}

	if resp.StatusCode != responseCode {
		t.Fatalf("expected %d, got %d\n%s", responseCode, resp.StatusCode, respBody)
	}

	cli.Cookies = append(cli.Cookies, resp.Cookies()...)

	return resp, respBody
}

// TestGet is a TestRequest shortcut for GET requests
func (cli *TestClient) TestGet(
	t testing.TB,
	endpoint string, responseCode int,
) (resp *http.Response, respBody []byte) {
	t.Helper()
	req, err := http.NewRequest(http.MethodGet, cli.testServer.URL+endpoint, nil)
	if err != nil {
		t.Fatal(err)
	}

	return cli.TestRequest(t, req, responseCode)
}

// TestPost is a TestRequest shortcut for POST requests
func (cli *TestClient) TestPost(
	t testing.TB,
	endpoint string, requestData interface{}, responseCode int,
) (resp *http.Response, respBody []byte) {
	t.Helper()

	body, err := json.Marshal(requestData)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodPost, cli.testServer.URL+endpoint, bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	return cli.TestRequest(t, req, responseCode)
}
