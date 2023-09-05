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
	var err error
	log := zerolog.New(os.Stdout)
	_ = os.Setenv("SERVICE_URL", ":8080")

	cfg := Config{}
	err = senv.Load(&cfg)
	if err != nil {
		log.Err(err).Send()
		//return exitErr
	}

	svc, err := newService(
		WithZerolog(log),
		WithServerURL(cfg.URL),
		//WithStorage("unix:///Users/goguson/.colima/default/docker.sock"))
		WithStorage(cfg.DockerHost))

	if err != nil {
		log.Err(err).Send()
		return exitErr
	}

	svc.run()

	err = svc.terminate()
	if err != nil {
		log.Error().Err(fmt.Errorf("svc.Terminate: %w", err))
	}
	if err != nil {
		return exitErr
	}
	return exitNoErr
}
