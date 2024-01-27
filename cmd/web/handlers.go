package main

import (
	"encoding/json"
	"errors"
	"example.com/practice-rest/internal/models"
	"example.com/practice-rest/pkg/lib"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

type PostDTO struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (app *application) getPosts(res http.ResponseWriter, _ *http.Request) {
	posts, err := app.post.Latest()

	if err != nil {
		app.errorLog.Println(err)
		lib.WriteJSON(res, http.StatusInternalServerError, lib.InternalServerError)
		return
	}

	lib.WriteJSON(res, http.StatusOK, lib.Response{Status: true, Result: posts, Message: "Posts Found"})
	return
}

func (app *application) getSinglePost(res http.ResponseWriter, req *http.Request) {
	// When httprouter is parsing a request, the values of any named parameters
	// will be stored in the request context. We'll talk about request context
	// in detail later in the book, but for now it's enough to know that you can
	// use the ParamsFromContext() function to retrieve a slice containing these
	// parameter names and values like so:
	params := httprouter.ParamsFromContext(req.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		app.errorLog.Println(err)
		lib.WriteJSON(res, http.StatusInternalServerError, lib.InternalServerError)
		return
	}

	post, err := app.post.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			lib.WriteJSON(res, http.StatusNotFound, lib.Response{Status: false, Result: nil, Message: "Post Not Found"})
			return
		}
		app.errorLog.Println(err)
		lib.WriteJSON(res, http.StatusInternalServerError, lib.InternalServerError)
		return

	}
	lib.WriteJSON(res, http.StatusOK, lib.Response{Status: true, Result: post, Message: "Post Found"})
}

func (app *application) createPost(res http.ResponseWriter, req *http.Request) {
	// Limit the size of the request body to 4KB
	req.Body = http.MaxBytesReader(res, req.Body, 4096)
	body := new(PostDTO)
	json.NewDecoder(req.Body).Decode(&body)

	fieldErrors := make(map[string]string)

	if body.Title == "" {
		fieldErrors["title"] = "Title cannot be blank"
	} else if len(body.Title) > 100 {
		fieldErrors["title"] = "This field is too long (maximum is 100 characters)"
	}

	if body.Content == "" {
		fieldErrors["content"] = "Content cannot be blank"
	}

	if len(fieldErrors) > 0 {
		lib.WriteJSON(res, http.StatusBadRequest, lib.Response{Status: false, Result: fieldErrors, Message: "Validation Error"})
		return
	}

	id, err := app.post.Insert(body.Title, body.Content)
	if err != nil {
		app.errorLog.Println(err)
		lib.WriteJSON(res, http.StatusInternalServerError, lib.InternalServerError)
		return
	}

	app.infoLog.Println(body)
	lib.WriteJSON(res, http.StatusOK, lib.Response{Status: true, Result: id, Message: "New Post Created"})
}

func healthCheck(res http.ResponseWriter, req *http.Request) {
	lib.WriteJSON(res, http.StatusOK, lib.Response{Status: true, Result: "Healthy", Message: "Hello World"})
	return
}
