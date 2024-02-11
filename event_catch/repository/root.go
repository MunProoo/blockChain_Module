package repository

import (
	"context"
	"event/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository struct {
	config *config.Config

	client *mongo.Client
	db     *mongo.Database

	Tx            *mongo.Collection
	NFT           *mongo.Collection
	NFTCollection *mongo.Collection
}

func NewRepository(config *config.Config) (*Repository, error) {
	r := Repository{
		config: config,
	}

	var err error
	ctx := context.Background()

	mongoConf := config.Database.Mongo

	if r.client, err = mongo.Connect(ctx, options.Client().ApplyURI(mongoConf.DataSource)); err != nil {
		return nil, err
	} else if err = r.client.Ping(ctx, nil); err != nil {
		return nil, err
	} else {
		r.db = r.client.Database(mongoConf.DB)

		r.Tx = r.db.Collection(mongoConf.Tx)
		r.NFT = r.db.Collection(mongoConf.NFT)
		r.NFTCollection = r.db.Collection(mongoConf.NFTCollection)
	}

	return &r, nil
}
