package config

import (
	"os"

	"github.com/naoina/toml"
)

type Config struct {
	Database struct {
		Mongo struct {
			DataSource string
			DB         string

			NFTCollection string
			NFT           string
			Tx            string
		}
	}

	Node struct {
		Uri        string
		StartBlock int64
	}
}

func NewConfig(path string) *Config {
	c := new(Config)
	if file, err := os.Open(path); err != nil {
		panic(err)
	} else if err = toml.NewDecoder(file).Decode(&c); err != nil {
		panic(err)
	} else {
		return c
	}

}
