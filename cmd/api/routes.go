package main

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()

	mux.NotFound(app.notFound)
	mux.MethodNotAllowed(app.methodNotAllowed)

	mux.Use(app.recoverPanic)
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, 
	}))
	mux.Get("/status", app.status)

	// User routes
	mux.Post("/users", app.registerUser)
	mux.Post("/users/login", app.loginUser)
	mux.With(app.authenticateToken).Get("/users/profile", app.getUserProfile)

	// Photo routes
	mux.With(app.authenticateToken).Post("/photos", app.uploadPhoto)
	mux.Get("/photos", app.getAllPhotos)
	mux.Get("/photos/{id}", app.getPhoto)
	mux.With(app.authenticateToken).Get("/users/photos", app.getUserPhotos)
	mux.Get("/photos/search", app.searchPhotos)

	// Token routes
	mux.Post("/tokens", app.createAuthenticationTokenHandler)

	// Like routes
	mux.With(app.authenticateToken).Post("/photos/{id}/like", app.likePhoto)
	mux.With(app.authenticateToken).Delete("/photos/{id}/like", app.unlikePhoto)

	// Comment routes
	mux.With(app.authenticateToken).Post("/photos/{id}/comments", app.addComment)
	mux.Get("/photos/{id}/comments", app.getPhotoComments)


	// Serve uploaded photos
	uploadDir := os.Getenv("UPLOAD_DIR")
	if uploadDir == "" {
		uploadDir = "/uploads"
	}
	fileServer := http.FileServer(http.Dir(uploadDir))
	mux.Handle("/uploads/*", http.StripPrefix("/uploads/", fileServer))

	return mux
}
