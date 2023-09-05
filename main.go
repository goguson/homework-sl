package main

import (
	"fmt"
	"github.com/goguson/senv"
	"github.com/rs/zerolog"
	"os"
)

var (
	exitNoErr = 0
	exitErr   = 1
)

func main() {
	os.Exit(Start())
}

func Start() int {
	log := zerolog.New(os.Stdout)

	cfg := Config{}
	err := senv.Load(&cfg)
	if err != nil {
		log.Err(err).Send()
	}

	svc, err := newService(
		WithZerolog(log),
		WithServerURL(cfg.URL),
		WithStorage(cfg.DockerHost))

	if err != nil {
		log.Err(err).Send()
		return exitErr
	}

	svc.run()

	err = svc.terminate()
	if err != nil {
		log.Err(fmt.Errorf("svc.Terminate: %w", err)).Send()
		return exitErr
	}
	return exitNoErr
}
