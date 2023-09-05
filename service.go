package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/goguson/homework-object-storage/storage"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Config struct {
	URL        string `senv:"SERVICE_URL"`
	DockerHost string `senv:"DOCKER_HOST"`
}

type Service struct {
	logger    zerolog.Logger
	store     storage.Store
	serverURL string
}

func (s *Service) run() {
	done := make(chan error)
	go func() {
		defer close(done)

		s.logger.Info().Msgf("url: %s", s.serverURL)

		err := http.ListenAndServe(s.serverURL, Routes(s))
		if err != nil {
			done <- err
			return
		}
		return
	}()

	exit := make(chan os.Signal, 1)
	var err error
	var sig os.Signal

	signal.Notify(exit, os.Interrupt, os.Kill, syscall.SIGTERM)
	select {
	case err = <-done:
		if err != nil {
			log.Error().Err(fmt.Errorf("svc.Run: %w", err))
		}
	case sig = <-exit:
		log.Info().Msgf("received stop signal: %s", sig)
	}
}

func newService(args ...func(*Service) error) (*Service, error) {
	svc := Service{}
	for _, arg := range args {
		if err := arg(&svc); err != nil {
			return nil, fmt.Errorf("args: %w", err)
		}
	}

	return &svc, nil
}

func (s *Service) terminate() error {
	s.logger.Warn().Msgf("%s graceful shutdown")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	stopSignal := make(chan struct{})
	go func() {
		defer close(stopSignal)
	}()

	select {
	case <-stopSignal:
	case <-ctx.Done():
		return fmt.Errorf("terminate: %w", ctx.Err())
	}
	return nil
}

func WithServerURL(url string) func(*Service) error {
	return func(svc *Service) error {
		if url == "" {
			return errors.New("server url is empty")
		}
		svc.serverURL = url
		return nil
	}
}

func WithStorage(host string) func(*Service) error {
	return func(svc *Service) error {
		c, err := storage.NewSocketClient(host)
		if err != nil {
			return fmt.Errorf("storage.NewSocketClient: %w", err)
		}

		svc.store = storage.NewStore(c)
		return nil
	}
}

func WithZerolog(l zerolog.Logger) func(*Service) error {
	return func(svc *Service) error {
		svc.logger = l
		return nil
	}
}
