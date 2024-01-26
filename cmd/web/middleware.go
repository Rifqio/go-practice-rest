package main

import "net/http"

// Add the secure headers middleware based on the OWASP specification
// https://owasp.org/www-project-secure-headers/index.html#configuration-proposal
func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self';")
		res.Header().Set("X-Content-Type-Options", "nosniff")
		res.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		res.Header().Set("X-Frame-Options", "deny")
		res.Header().Set("X-XSS-Protection", "1; mode=block")

		next.ServeHTTP(res, req)
	})
}

func (app *application) requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		app.httpLog.Printf("%s - %s %s %s", req.RemoteAddr, req.Proto, req.Method, req.URL.RequestURI())
		next.ServeHTTP(res, req)
	})
}
