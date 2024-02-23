package handlers

import (
	"github.com/go-chi/chi/v5"
)

func (uh *URLHandler) InitRouter() chi.Router {

	router := chi.NewRouter()
	// router.Use(middleware.Logger)
	router.Use(uh.Logger)
	router.Post("/", uh.GenerateShortener)
	router.Get("/{id}", uh.GetEndpointByShortener)
	router.Post("/shorten", uh.GenerateJSONShortener)

	return router
}
