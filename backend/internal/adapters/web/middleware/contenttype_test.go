package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestContentTypeValidation_JSONAccepted(t *testing.T) {
	handler := ContentTypeValidation(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodPost, "/graphql", strings.NewReader(`{"query":"{ __typename }"}`))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

func TestContentTypeValidation_JSONWithCharsetAccepted(t *testing.T) {
	handler := ContentTypeValidation(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodPost, "/graphql", strings.NewReader(`{"query":"{ __typename }"}`))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

func TestContentTypeValidation_FormEncodedRejected(t *testing.T) {
	handler := ContentTypeValidation(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not be called for form-encoded POST")
	}))

	req := httptest.NewRequest(http.MethodPost, "/graphql", strings.NewReader("query=%7B+__typename+%7D"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnsupportedMediaType {
		t.Errorf("expected 415, got %d", rr.Code)
	}
	if !strings.Contains(rr.Body.String(), "Content-Type must be application/json") {
		t.Errorf("expected JSON error message, got %s", rr.Body.String())
	}
}

func TestContentTypeValidation_MissingContentTypeRejected(t *testing.T) {
	handler := ContentTypeValidation(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not be called for POST without Content-Type")
	}))

	req := httptest.NewRequest(http.MethodPost, "/graphql", strings.NewReader(`{"query":"{ __typename }"}`))
	// No Content-Type header set
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnsupportedMediaType {
		t.Errorf("expected 415, got %d", rr.Code)
	}
}

func TestContentTypeValidation_GETPassesThrough(t *testing.T) {
	handler := ContentTypeValidation(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/graphql", nil)
	// No Content-Type needed for GET
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}
