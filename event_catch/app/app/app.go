package app

import (
	"event/config"
	"event/event"
	"event/repository"

	ethTypes "github.com/hacpy/go-ethereum/core/types"
	"github.com/hacpy/go-ethereum/ethclient"
)

type App struct {
	config *config.Config

	client *ethclient.Client

	repository *repository.Repository
	scan       *event.Scan
	catch      *event.Catch
}

func NewApp(config *config.Config) *App {
	a := App{
		config: config,
	}

	var err error

	// infura mumbai 연결
	if a.client, err = ethclient.Dial(config.Node.Uri); err != nil {
		panic(err)
	}
	// DB 연결
	if a.repository, err = repository.NewRepository(config); err != nil {
		panic(err)
	}

	var eventChan chan []ethTypes.Log
	// 캐치 연결 (어떤 이벤트를 캐치할 건지를 정하고 나서 스캐너 연결하므로 먼저 연결)
	if a.catch, err = event.NewCatch(config, a.client, eventChan); err != nil {
		panic(err)
	}
	// 스캐너 연결
	if a.scan, eventChan, err = event.NewScan(config, a.client, a.catch.GetEventToCatch()); err != nil {
		panic(err)
	}

	return &a
}
