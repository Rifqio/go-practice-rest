package main

import (
	"encoding/json"
	"errors"
	"example.com/practice-rest/internal/models"
	"example.com/practice-rest/internal/validator"
	"example.com/practice-rest/pkg/lib"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/samber/lo"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
)

func (app *application) getPosts(res http.ResponseWriter, _ *http.Request) {
	posts, err := app.post.Latest()

	if lo.IsNotEmpty(err) {
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
	if lo.IsNotEmpty(err) {
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
	type PostDTO struct {
		Title   string `json:"title" validate:"required"`
		Content string `json:"content" validate:"required"`
		validator.Validator
	}

	// Limit the size of the request body to 4KB
	req.Body = http.MaxBytesReader(res, req.Body, 4096)
	body := new(PostDTO)
	json.NewDecoder(req.Body).Decode(&body)

	body.CheckField(validator.NotEmpty(body.Title), "title", "title cannot be blank")
	body.CheckField(validator.MaxChars(body.Title, 100), "title", "this field is too long (maximum is 100 characters)")

	body.CheckField(validator.NotEmpty(body.Content), "content", "content cannot be blank")

	if !body.Valid() {
		lib.WriteJSON(res, http.StatusBadRequest, lib.Response{Status: false, Result: body.Errors, Message: "Validation Error"})
		return
	}

	id, err := app.post.Insert(body.Title, body.Content)
	if lo.IsNotEmpty(err) {
		app.errorLog.Println(err)
		lib.WriteJSON(res, http.StatusInternalServerError, lib.InternalServerError)
		return
	}

	app.sessionManger.Put(req.Context(), "flash", "Post Created Successfully")

	app.infoLog.Println(body)
	lib.WriteJSON(res, http.StatusOK, lib.Response{Status: true, Result: id, Message: "New Post Created"})
}

func (app *application) userSignup(res http.ResponseWriter, req *http.Request) {
	type UserSignupDTO struct {
		Name     string `json:"name" validate:"required"`
		Email    string `json:"email" validate:"required"`
		Password string `json:"password" validate:"required"`
		validator.Validator
	}

	body := new(UserSignupDTO)
	json.NewDecoder(req.Body).Decode(&body)

	body.CheckField(validator.NotEmpty(body.Name), "name", "name cannot be blank")
	body.CheckField(validator.MaxChars(body.Name, 100), "name", "this field is too long (maximum is 100 characters)")

	body.CheckField(validator.NotEmpty(body.Email), "email", "email cannot be blank")
	body.CheckField(validator.Matches(body.Email, validator.EmailRX), "email", "email is not valid")

	body.CheckField(validator.NotEmpty(body.Password), "password", "password cannot be blank")
	body.CheckField(validator.MinChars(body.Password, 8), "password", "password must be at least 8 characters")

	if !body.Valid() {
		lib.WriteJSON(res, http.StatusBadRequest, lib.Response{Status: false, Result: body.Errors, Message: "Validation Error"})
		return
	}

	hashedPass, errHashing := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	id, errInserting := app.user.Insert(body.Name, body.Email, string(hashedPass))

	if err := errors.Join(errHashing, errInserting); lo.IsNotEmpty(err) {
		app.errorLog.Println(err)
		lib.WriteJSON(res, http.StatusInternalServerError, lib.InternalServerError)
		return
	}

	app.infoLog.Println("User Signup Successfully", id, body.Email)
	lib.WriteJSON(res, http.StatusOK, lib.Response{Status: true, Result: body.Email, Message: "User Signup Successfully"})
}

func (app *application) userLogin(res http.ResponseWriter, req *http.Request) {
	fmt.Println("User Login")
}

func (app *application) userLogout(res http.ResponseWriter, req *http.Request) {
	fmt.Println("User Logout")
}

func healthCheck(res http.ResponseWriter, req *http.Request) {
	lib.WriteJSON(res, http.StatusOK, lib.Response{Status: true, Result: "Healthy", Message: "Hello World"})
	return
}
