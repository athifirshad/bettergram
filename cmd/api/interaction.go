package main

import (
	"errors"
	"net/http"

	"athifirshad.com/bettergram/internal/data"
	"athifirshad.com/bettergram/internal/request"
	"athifirshad.com/bettergram/internal/response"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (app *application) likePhoto(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)
	if user == data.AnonymousUser {
		app.invalidAuthenticationToken(w, r)
		return
	}

	photoID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		app.notFound(w, r)
		return
	}

	like := &data.Like{
		PhotoID: photoID,
		UserID:  user.ID,
	}

	err = app.data.Likes.Insert(like)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateLike):
			app.errorMessage(w, r, http.StatusConflict, "post has already been liked", nil)
		default:
			app.serverError(w, r, err)
		}
		return
	}

	err = response.JSON(w, http.StatusCreated, like)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) unlikePhoto(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)
	if user == data.AnonymousUser {
		app.invalidAuthenticationToken(w, r)
		return
	}

	photoID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		app.notFound(w, r)
		return
	}

	err = app.data.Likes.Delete(photoID, user.ID)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (app *application) addComment(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)
	if user == data.AnonymousUser {
		app.invalidAuthenticationToken(w, r)
		return
	}

	photoID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		app.notFound(w, r)
		return
	}

	var input struct {
		Content string `json:"content"`
	}

	err = request.DecodeJSON(w, r, &input)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	comment := &data.Comment{
		PhotoID: photoID,
		UserID:  user.ID,
		Content: input.Content,
	}

	err = app.data.Comments.Insert(comment)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	err = response.JSON(w, http.StatusCreated, comment)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) getPhotoComments(w http.ResponseWriter, r *http.Request) {
	photoID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		app.notFound(w, r)
		return
	}

	comments, err := app.data.Comments.GetByPhotoID(photoID)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			app.notFound(w, r)
		default:
			app.serverError(w, r, err)
		}
		return
	}

	err = response.JSON(w, http.StatusOK, comments)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) searchPhotos(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		app.badRequest(w, r, errors.New("search query is required"))
		return
	}

	photos, err := app.data.Photos.Search(query)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	err = response.JSON(w, http.StatusOK, photos)
	if err != nil {
		app.serverError(w, r, err)
	}
}
