package main

import (
	"errors"
	"net/http"
	"time"

	"athifirshad.com/bettergram/internal/data"
	"athifirshad.com/bettergram/internal/request"
	"athifirshad.com/bettergram/internal/response"
)

func (app *application) createAuthenticationTokenHandler(w http.ResponseWriter, r *http.Request) {
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
		return
	}

	err = response.JSON(w, http.StatusCreated, map[string]interface{}{"authentication_token": token})
	if err != nil {
		app.serverError(w, r, err)
	}
}