package req_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/digiconvent/req"
)

func TestPost_SendsMethodAndBody(t *testing.T) {
	var gotMethod, gotBody string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		body, _ := io.ReadAll(r.Body)
		gotBody = string(body)
	}))
	defer srv.Close()

	_, _, err := req.Post(srv.URL, strings.NewReader("hello")).Do()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if gotMethod != http.MethodPost {
		t.Errorf("expected method POST, got %q", gotMethod)
	}
	if gotBody != "hello" {
		t.Errorf("expected body %q, got %q", "hello", gotBody)
	}
}

func TestPost_NonSuccessStatus_ReturnsErrorMentioningPost(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("nope"))
	}))
	defer srv.Close()

	_, _, err := req.Post(srv.URL, strings.NewReader("x")).Do()
	if err == nil {
		t.Fatal("expected an error for a 400 response")
	}
	if !strings.Contains(err.Error(), "POST") || !strings.Contains(err.Error(), "nope") {
		t.Errorf("expected the error to mention POST and the body, got: %v", err)
	}
}
