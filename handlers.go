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
			log := logger.With().
				Timestamp().
				Str("method", r.Method).
				Str("path", r.URL.Path).
				Logger()
			ctx := log.WithContext(r.Context())
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func Routes(svc *Service) *chi.Mux {
	r := chi.NewRouter()
	r.Use(ZerologMiddleware(svc.logger))
	r.Use(middleware.RedirectSlashes)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))

	// Define the route for handling PUT requests
	r.Put("/object/{id:[a-zA-Z0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		objectID := chi.URLParam(r, "id")

		err := svc.store.Put(ctx, objectID, r)
		if err != nil {
			zerolog.Ctx(r.Context()).Error().Err(err).Msg("error putting object")
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
			zerolog.Ctx(r.Context()).Error().Err(err).Send()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	r.Get("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		// if there was a database connection, any dependency with state and so on, we would check it here
		w.WriteHeader(http.StatusOK)
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		logger := zerolog.Ctx(r.Context())
		logger.Info().Msg("hello world")
		w.WriteHeader(http.StatusOK)
	})

	return r
}
