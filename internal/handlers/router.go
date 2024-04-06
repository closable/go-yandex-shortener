package handlers

import (
	"github.com/go-chi/chi/v5"
)

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

	// router.Group(func(r chi.Router) {
	// 	r.Use(uh.Auth)
	// 	r.Get("/api/user/urls", uh.GetUrls)
	// 	r.Delete("/api/user/urls", uh.DelUrls)
	// })

	return router
}
