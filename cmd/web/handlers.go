package main

import (
	"encoding/json"
	"errors"
	"example.com/practice-rest/internal/models"
	"example.com/practice-rest/pkg/lib"
	"net/http"
	"strconv"
)

type PostDTO struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (app *application) getPost(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		app.errorLog.Println("Method Not Allowed")
		lib.WriteJSON(res, http.StatusMethodNotAllowed, lib.MethodNotAllowed)
		return
	}

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
	if req.Method != http.MethodGet {
		app.errorLog.Println("Method Not Allowed")
		lib.WriteJSON(res, http.StatusMethodNotAllowed, lib.MethodNotAllowed)
		return
	}

	id, err := strconv.Atoi(req.URL.Query().Get("id"))
	if err != nil || id < 1 {
		lib.WriteJSON(res, http.StatusBadRequest, lib.Response{Status: false, Result: nil, Message: "Invalid Post ID"})
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
	if req.Method != http.MethodPost {
		app.errorLog.Println("Method Not Allowed")
		lib.WriteJSON(res, http.StatusMethodNotAllowed, lib.MethodNotAllowed)
		return
	}

	body := new(PostDTO)
	json.NewDecoder(req.Body).Decode(&body)

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
