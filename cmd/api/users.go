package main

import (
	"errors"
	"net/http"
	"time"

	"athifirshad.com/bettergram/internal/data"
	"athifirshad.com/bettergram/internal/request"
	"athifirshad.com/bettergram/internal/response"
)

type UserController struct {
	UserModel *data.UserModel
}

func (app *application) registerUser(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := request.DecodeJSON(w, r, &input)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	user := &data.User{
		Username: input.Username,
		Email:    input.Email,
	}

	err = user.Password.Set(input.Password)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	err = app.data.Users.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			app.badRequest(w, r, errors.New("a user with this email address already exists"))
		case errors.Is(err, data.ErrDuplicateUsername):
			app.badRequest(w, r, errors.New("a user with this username already exists"))
		default:
			app.serverError(w, r, err)
		}
		return
	}

	err = response.JSON(w, http.StatusCreated, user)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {
	
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := request.DecodeJSON(w, r, &input)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	user, err := app.data.Users.GetByEmail(input.Email)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			response.InvalidCredentials(w, r)
		default:
			app.serverError(w, r, err)
		}
		return
	}

	match, err := user.Password.Matches(input.Password)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	if !match {
		response.InvalidCredentials(w, r)
		return
	}

	token, err := app.data.Tokens.New(user.ID, 24*time.Hour, data.ScopeAuthentication)
	if err != nil {
		app.serverError(w, r, err)
	}

	err = response.JSON(w, http.StatusOK, map[string]string{"authentication_token": token.Plaintext})
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) getUserProfile(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)
	err := response.JSON(w, http.StatusOK, user)
	if err != nil {
		app.serverError(w, r, err)
	}
}
