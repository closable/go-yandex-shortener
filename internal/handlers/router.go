package handlers

import "github.com/go-chi/chi/v5"

func (h *URLHandler) InitRouter() chi.Router {

	router := chi.NewRouter()

	router.Post("/", h.GenerateShortener)
	router.Get("/{id}", h.GetEndpointByShortener)

	return router
}
