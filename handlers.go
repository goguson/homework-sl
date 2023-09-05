package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/minio/minio-go/v7"
	"github.com/rs/zerolog"
	"io"
	"net/http"
	"time"
)

const uploadErrMsg = "could not read form file"
const downloadErrMsg = "problem occurred during download of data from another server"
const fileReadErrMsg = "could not read file from storage"

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

	r.Put("/object/{id:[a-zA-Z0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := zerolog.Ctx(ctx)
		objectID := chi.URLParam(r, "id")

		file, header, err := r.FormFile("file")
		if err != nil {
			logger.Err(err).Send()
			http.Error(w, fileReadErrMsg, http.StatusInternalServerError)
		}
		defer file.Close()

		err = svc.store.Put(ctx, objectID, file, header)
		if err != nil {
			logger.Err(err).Send()
			http.Error(w, uploadErrMsg, http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	})

	r.Get("/object/{id:[a-zA-Z0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		objectID := chi.URLParam(r, "id")
		logger := zerolog.Ctx(ctx)

		reader, err := svc.store.Get(ctx, objectID)
		defer reader.Close()

		stat, err := reader.Stat()
		if err != nil {
			logger.Err(fmt.Errorf("reader.Stat: %w", err)).Send()
			mErr := minio.ToErrorResponse(err)
			if mErr.Code == "NoSuchKey" {
				http.Error(w, "object not found", http.StatusNotFound)
				return
			}
			http.Error(w, downloadErrMsg, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", stat.ContentType)
		_, err = io.Copy(w, reader)
		if err != nil {
			logger.Err(fmt.Errorf("io.Copy: %w", err)).Send()
			http.Error(w, downloadErrMsg, http.StatusInternalServerError)
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
		_, _ = w.Write([]byte("hello world"))
	})

	return r
}
