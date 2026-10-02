package req_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/digiconvent/req"
)

func TestDo_SuccessReturnsBody(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
	}))
	defer srv.Close()

	body, _, err := req.Get(srv.URL).Do()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if string(body) != "hello" {
		t.Errorf("expected body %q, got %q", "hello", string(body))
	}
}

func TestDo_SendsChainedHeaders(t *testing.T) {
	var gotAuth string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
	}))
	defer srv.Close()

	_, _, err := req.Get(srv.URL).Header("Authorization", "Bearer secret").Do()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if gotAuth != "Bearer secret" {
		t.Errorf("expected the Authorization header to reach the server, got %q", gotAuth)
	}
}

func TestDo_MultipleChainedHeaders_AllReachTheServer(t *testing.T) {
	var gotAccept, gotAuth string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAccept = r.Header.Get("Accept")
		gotAuth = r.Header.Get("Authorization")
	}))
	defer srv.Close()

	_, _, err := req.Get(srv.URL).Header("Accept", "application/vnd.github.v3+json").Header("Authorization", "Bearer secret").Do()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if gotAccept != "application/vnd.github.v3+json" || gotAuth != "Bearer secret" {
		t.Errorf("expected both headers to reach the server, got accept=%q auth=%q", gotAccept, gotAuth)
	}
}

func TestDo_NonSuccessStatus_ReturnsStatusErrorWithBody(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("not here"))
	}))
	defer srv.Close()

	_, _, err := req.Get(srv.URL).Do()
	if err == nil {
		t.Fatal("expected an error for a 404 response")
	}

	var reqErr *req.RequestError
	if !errors.As(err, &reqErr) {
		t.Fatalf("expected a *req.RequestError, got: %T", err)
	}
	if reqErr.Stage != req.StageStatus {
		t.Errorf("expected stage %q, got %q", req.StageStatus, reqErr.Stage)
	}

	var statusErr *req.StatusError
	if !errors.As(err, &statusErr) {
		t.Fatalf("expected a *req.StatusError, got: %T", err)
	}
	if statusErr.StatusCode != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", statusErr.StatusCode)
	}
	if string(statusErr.Body) != "not here" {
		t.Errorf("expected body %q, got %q", "not here", string(statusErr.Body))
	}
}

func TestDo_InvalidURL_ReturnsBuildRequestStage(t *testing.T) {
	_, _, err := req.Get("://not-a-url").Do()
	if err == nil {
		t.Fatal("expected an error for a malformed URL")
	}

	var reqErr *req.RequestError
	if !errors.As(err, &reqErr) {
		t.Fatalf("expected a *req.RequestError, got: %T", err)
	}
	if reqErr.Stage != req.StageBuildRequest {
		t.Errorf("expected stage %q, got %q", req.StageBuildRequest, reqErr.Stage)
	}
}

func TestDo_UnreachableHost_ReturnsSendRequestStage(t *testing.T) {
	_, _, err := req.Get("http://127.0.0.1:1").Do()
	if err == nil {
		t.Fatal("expected an error for an unreachable host")
	}

	var reqErr *req.RequestError
	if !errors.As(err, &reqErr) {
		t.Fatalf("expected a *req.RequestError, got: %T", err)
	}
	if reqErr.Stage != req.StageSendRequest {
		t.Errorf("expected stage %q, got %q", req.StageSendRequest, reqErr.Stage)
	}
}
