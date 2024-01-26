package main

import (
	"example.com/practice-rest/pkg/lib"
	"fmt"
	"net/http"
	"runtime/debug"
)

func (app *application) serverError(res http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	lib.WriteJSON(res, http.StatusInternalServerError, lib.Response{Status: false, Result: nil, Message: "Internal Server Error"})
}

func (app *application) pageNotFound(res http.ResponseWriter) {
	lib.WriteJSON(res, http.StatusNotFound, lib.Response{Status: false, Result: nil, Message: "Page Not Found"})
}
