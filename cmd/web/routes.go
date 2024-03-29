package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"net/http"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	dynamic := alice.New(app.sessionManger.LoadAndSave)
	router.HandlerFunc(http.MethodGet,"/", healthCheck)
	router.Handler(http.MethodGet,"/post/:id", dynamic.ThenFunc(app.getSinglePost))
	router.Handler(http.MethodGet,"/post", dynamic.ThenFunc(app.getPosts))
	router.Handler(http.MethodPost, "/post", dynamic.ThenFunc(app.createPost))

	router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(app.userSignup))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.userLogin))
	router.Handler(http.MethodPost, "/user/logout", dynamic.ThenFunc(app.userLogout))

	router.NotFound = http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		app.pageNotFound(res)
	})

	standard := alice.New(app.recoverPanic, app.requestLogger, secureHeaders)
	return standard.Then(router)
}
