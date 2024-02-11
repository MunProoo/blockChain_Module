package event

import (
	"event/config"

	"github.com/hacpy/go-ethereum/ethclient"
)

type Scan struct {
	config *config.Config

	client *ethclient.Client
}

func NewScan(config *config.Config, client *ethclient.Client) (*Scan, error) {
	s := &Scan{
		config: config,
		client: client,
	}

	var err error

	if s.client, err = ethclient.Dial(config.Node.Uri); err != nil {
		return nil, err
	} else {
		return s, nil
	}

}
