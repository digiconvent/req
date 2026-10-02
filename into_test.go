package req_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/digiconvent/req"
)

func TestInto_SuccessUnmarshalsBody(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"name":"frodo"}`))
	}))
	defer srv.Close()

	var got struct {
		Name string `json:"name"`
	}
	if err := req.Get(srv.URL).Into(&got); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Name != "frodo" {
		t.Errorf("expected name %q, got %q", "frodo", got.Name)
	}
}

func TestInto_NonSuccessStatus_ReturnsStatusError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	var got struct{}
	err := req.Get(srv.URL).Into(&got)
	if err == nil {
		t.Fatal("expected an error for a 404 response")
	}
}
