package app

import (
	"event/config"
	"event/event"
	"event/repository"

	"github.com/hacpy/go-ethereum/ethclient"
)

type App struct {
	config *config.Config

	client *ethclient.Client

	repository *repository.Repository
	scan       *event.Scan
}

func NewApp(config *config.Config) *App {
	a := App{
		config: config,
	}

	var err error

	if a.repository, err = repository.NewRepository(config); err != nil {
		panic(err)
	}

	if a.client, err = ethclient.Dial(config.Node.Uri); err != nil {
		panic(err)
	}
	if a.scan, err = event.NewScan(config, a.client); err != nil {
		panic(err)
	}

	return &a
}
