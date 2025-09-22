package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type Response = httptest.ResponseRecorder

func newRequest(method string, url string, data string) *Response {
	var body *io.Reader

	if len(strings.TrimSpace(data)) != 0 {
		body = strings.NewReader(data)
	}

	req, err := http.NewRequest(method, url, bo)

}

func makeTestRequest(req *http.Request) *Response {

}

func TestHealthHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/health/", nil)

	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	Health(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d but got %d",
			http.StatusOK,
			w.Code)
	}
}

func TestRequestHandler(t *testing.T) {
	expected := "whats up"

	body := strings.NewReader("whats up")
	req, err := http.NewRequest("POST", "/", body)

	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()

	RequestHandler(w, req)

	res := w.Result()
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)

	if err != nil {
		t.Errorf("Error: %v", err)
	}

	actual := string(data)

	if actual != expected {
		t.Errorf("Expected %#v actual %#v", expected, actual)
	}
}
