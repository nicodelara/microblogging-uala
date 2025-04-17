package mongo

import (
	"context"

	"github.com/nicodelara/uala-challenge/internal/tweets/domain"
	"github.com/nicodelara/uala-challenge/internal/tweets/domain/ports"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoTweetRepository struct {
	collection *mongo.Collection
}

func NewMongoTweetRepository(client *mongo.Client, dbName, collName string) (ports.TweetRepository, error) {
	db := client.Database(dbName)
	collection := db.Collection(collName)

	// Crear índices
	indexModels := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "username", Value: 1},
				{Key: "createdAt", Value: -1},
			},
		},
	}

	_, err := collection.Indexes().CreateMany(context.Background(), indexModels)
	if err != nil {
		return nil, err
	}

	return &mongoTweetRepository{
		collection: collection,
	}, nil
}

func (r *mongoTweetRepository) SaveTweet(ctx context.Context, tweet *domain.Tweet) error {
	_, err := r.collection.InsertOne(ctx, tweet)
	return err
}
