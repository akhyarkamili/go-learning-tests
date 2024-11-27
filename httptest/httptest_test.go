package httptest

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHttpTest(t *testing.T) {
	// Act
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, world"))
	})
	server := httptest.NewServer(handler)
	// Assert
	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", resp.Status)
	}
	if resp.Body == nil {
		t.Error("response body is nil")
	}
	t.Log("TestHttpTest")
}
