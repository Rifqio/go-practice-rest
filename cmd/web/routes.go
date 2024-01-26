package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", healthCheck)
	mux.HandleFunc("/post", app.getSinglePost)
	mux.HandleFunc("/posts", app.getPost)
	mux.HandleFunc("/post/create", app.createPost)

	return mux
}
