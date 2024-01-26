package main

import "net/http"

// Update the signature for the routes() method so that it returns a
// http.Handler instead of *http.ServeMux.
func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", healthCheck)
	mux.HandleFunc("/post", app.getSinglePost)
	mux.HandleFunc("/posts", app.getPost)
	mux.HandleFunc("/post/create", app.createPost)

	return app.requestLogger(secureHeaders(mux))
}
