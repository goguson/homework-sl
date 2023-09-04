package main

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rs/zerolog"
	"net/http"
	"time"
)

func ZerologMiddleware(logger zerolog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r.WithContext(logger.WithContext(r.Context())))
		})
	}
}

func Routes(svc *Service) *chi.Mux {
	router := chi.NewRouter()
	//router.Use(ZerologMiddleware(logger))
	router.Use(middleware.RedirectSlashes)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(5 * time.Second))
	r := chi.NewRouter()

	// Define the route for handling PUT requests
	r.Put("/object/{id:[a-zA-Z0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		objectID := chi.URLParam(r, "id")

		err := svc.store.Put(ctx, objectID, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	})

	// Define the route for handling GET requests
	r.Get("/object/{id:[a-zA-Z0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		objectID := chi.URLParam(r, "id")

		err := svc.store.Get(ctx, objectID, w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNotImplemented)
	})

	return router
}
