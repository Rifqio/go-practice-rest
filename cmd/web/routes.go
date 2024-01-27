package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"net/http"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()
	router.HandlerFunc(http.MethodGet,"/", healthCheck)
	router.HandlerFunc(http.MethodGet,"/post/:id", app.getSinglePost)
	router.HandlerFunc(http.MethodGet,"/post", app.getPosts)
	router.HandlerFunc(http.MethodPost, "/post", app.createPost)

	router.NotFound = http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		app.pageNotFound(res)
	})
	standard := alice.New(app.recoverPanic, app.requestLogger, secureHeaders)
	return standard.Then(router)
}
