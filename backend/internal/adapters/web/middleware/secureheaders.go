package middleware

import (
	"net/http"
	"os"

	"github.com/unrolled/secure"
)

// SecureHeaders returns middleware setting security headers
// (HSTS, X-Content-Type-Options, X-Frame-Options, CSP, XSS filter).
// In development mode, HSTS is disabled to allow localhost.
func SecureHeaders() func(http.Handler) http.Handler {
	isDevelopment := os.Getenv("APP_ENV") != "production"

	secureMiddleware := secure.New(secure.Options{
		// HSTS: enforce HTTPS for 1 year, include subdomains (only in production)
		STSSeconds:           31536000,
		STSIncludeSubdomains: true,
		STSPreload:           true,

		// Prevent MIME sniffing
		ContentTypeNosniff: true,

		// Prevent clickjacking
		FrameDeny: true,

		// XSS protection (legacy browsers)
		BrowserXssFilter: true,

		// Content Security Policy for backend responses
		ContentSecurityPolicy: "default-src 'self'; script-src 'self'; style-src 'self'; img-src 'self' https://i.ytimg.com https://yt3.ggpht.com; connect-src 'self'; font-src 'self';",

		// Detect HTTPS behind Sevalla/Cloudflare reverse proxy
		SSLProxyHeaders: map[string]string{"X-Forwarded-Proto": "https"},

		// Only enforce in production (development allows http://localhost)
		IsDevelopment: isDevelopment,
	})

	return secureMiddleware.Handler
}
