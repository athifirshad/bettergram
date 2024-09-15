package main

import (
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"athifirshad.com/bettergram/internal/data"
	"athifirshad.com/bettergram/internal/response"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type createPhotoInput struct {
	Caption string `json:"caption"`
}

func (app *application) uploadPhoto(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)
	if user == data.AnonymousUser {
		app.invalidAuthenticationToken(w, r)
		return
	}

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	file, handler, err := r.FormFile("photo")
	if err != nil {
		app.badRequest(w, r, err)
		return
	}
	defer file.Close()

	caption := r.FormValue("caption")

	uploadDir := os.Getenv("UPLOAD_DIR")
	if uploadDir == "" {
		uploadDir = "/uploads"
	}

	filename := uuid.New().String() + filepath.Ext(handler.Filename)
	filepath := filepath.Join(uploadDir, filename)

	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		app.serverError(w, r, err)
		return
	}

	out, err := os.Create(filepath)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	photo := &data.Photo{
		UserID:   user.ID,
		PhotoURL: "/uploads/" + filename,
		Caption:  caption,
	}

	err = app.data.Photos.Insert(photo)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	err = response.JSON(w, http.StatusCreated, photo)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) getPhoto(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	photoID, err := uuid.Parse(idStr)
	if err != nil {
		app.notFound(w, r)
		return
	}

	photo, err := app.data.Photos.GetByID(photoID)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			app.notFound(w, r)
		default:
			app.serverError(w, r, err)
		}
		return
	}

	err = response.JSON(w, http.StatusOK, photo)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) getUserPhotos(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)
	if user == data.AnonymousUser {
		app.invalidAuthenticationToken(w, r)
		return
	}

	photos, err := app.data.Photos.GetByUserID(user.ID)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	err = response.JSON(w, http.StatusOK, photos)
	if err != nil {
		app.serverError(w, r, err)
	}
}


func (app *application) getAllPhotos(w http.ResponseWriter, r *http.Request) {
	photos, err := app.data.Photos.GetAll()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	err = response.JSON(w, http.StatusOK, photos)
	if err != nil {
		app.serverError(w, r, err)
	}
}
