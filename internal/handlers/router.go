package handlers

import (
	"github.com/go-chi/chi/v5"
)

// InitRouter инициализация роутера
func (uh *URLHandler) InitRouter() chi.Router {

	router := chi.NewRouter()
	// router.Use(middleware.Logger)

	router.Use(uh.Compressor)
	router.Use(uh.Logger)
	router.Use(uh.Auth)

	router.Post("/", uh.GenerateShortener)
	router.Get("/{id}", uh.GetEndpointByShortener)
	router.Get("/ping", uh.CheckBaseActivity)
	router.Post("/api/shorten", uh.GenerateJSONShortener)
	router.Post("/api/shorten/batch", uh.UploadBatch)

	router.Get("/api/user/urls", uh.GetUrls)
	router.Delete("/api/user/urls", uh.DelUrls)

	router.Group(func(r chi.Router) {
		r.Use(uh.Traster)
		r.Get("/api/internal/stats", uh.GetStats)

	})

	return router
}
