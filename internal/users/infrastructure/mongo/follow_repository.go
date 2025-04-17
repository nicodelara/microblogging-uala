package mongo

import (
	"context"
	"errors"

	"github.com/nicodelara/microblogging-uala/internal/users/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoFollowRepository struct {
	client     *mongo.Client
	dbName     string
	collection string
}

func NewMongoFollowRepository(client *mongo.Client, dbName, collection string) (*mongoFollowRepository, error) {
	if client == nil {
		return nil, errors.New("mongo client is required")
	}
	if dbName == "" {
		return nil, errors.New("database name is required")
	}
	if collection == "" {
		return nil, errors.New("collection name is required")
	}

	// Crear índices
	coll := client.Database(dbName).Collection(collection)
	indexModels := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "username", Value: 1},
				{Key: "following", Value: 1},
			},
			Options: options.Index().SetUnique(true),
		},
	}

	_, err := coll.Indexes().CreateMany(context.Background(), indexModels)
	if err != nil {
		return nil, err
	}

	return &mongoFollowRepository{
		client:     client,
		dbName:     dbName,
		collection: collection,
	}, nil
}

func (r *mongoFollowRepository) GetFollowings(ctx context.Context, username string) ([]string, error) {
	collection := r.client.Database(r.dbName).Collection(r.collection)
	cursor, err := collection.Find(ctx, bson.M{"username": username})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var follows []domain.Follow
	if err := cursor.All(ctx, &follows); err != nil {
		return nil, err
	}

	followingUsernames := make([]string, len(follows))
	for i, follow := range follows {
		followingUsernames[i] = follow.Following
	}

	return followingUsernames, nil
}

func (r *mongoFollowRepository) FollowUser(ctx context.Context, follow *domain.Follow) error {
	collection := r.client.Database(r.dbName).Collection(r.collection)
	_, err := collection.InsertOne(ctx, follow)
	return err
}
