package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestSecureHeaders_ContentTypeNosniff(t *testing.T) {
	handler := SecureHeaders()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/graphql", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if got := rr.Header().Get("X-Content-Type-Options"); got != "nosniff" {
		t.Errorf("expected X-Content-Type-Options=nosniff, got %q", got)
	}
}

func TestSecureHeaders_FrameDeny(t *testing.T) {
	handler := SecureHeaders()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/graphql", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if got := rr.Header().Get("X-Frame-Options"); got != "DENY" {
		t.Errorf("expected X-Frame-Options=DENY, got %q", got)
	}
}

func TestSecureHeaders_ProductionHSTS(t *testing.T) {
	t.Setenv("APP_ENV", "production")

	// Must recreate middleware after env change (reads env at init)
	handler := SecureHeaders()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/graphql", nil)
	req.Header.Set("X-Forwarded-Proto", "https")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	sts := rr.Header().Get("Strict-Transport-Security")
	if sts == "" {
		t.Error("expected Strict-Transport-Security header in production, got empty")
	}
}

func TestSecureHeaders_DevelopmentNoHSTS(t *testing.T) {
	// Ensure we're in development mode
	originalEnv := os.Getenv("APP_ENV")
	t.Setenv("APP_ENV", "development")
	defer func() {
		if originalEnv != "" {
			os.Setenv("APP_ENV", originalEnv)
		}
	}()

	handler := SecureHeaders()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/graphql", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	// In development mode, unrolled/secure skips HSTS and other production-only headers
	// The IsDevelopment flag disables enforcement
	sts := rr.Header().Get("Strict-Transport-Security")
	if sts != "" {
		t.Errorf("expected no Strict-Transport-Security in development, got %q", sts)
	}
}

func TestSecureHeaders_CSPHeaderSet(t *testing.T) {
	handler := SecureHeaders()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/graphql", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	csp := rr.Header().Get("Content-Security-Policy")
	if csp == "" {
		t.Error("expected Content-Security-Policy header, got empty")
	}
}
