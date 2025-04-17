package mongo

import (
	"context"
	"errors"

	"github.com/nicodelara/uala-challenge/internal/users/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoUserRepository struct {
	client     *mongo.Client
	dbName     string
	collection string
}

func NewMongoUserRepository(client *mongo.Client, dbName, collection string) (*mongoUserRepository, error) {
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
			Keys:    bson.D{{Key: "username", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	}

	_, err := coll.Indexes().CreateMany(context.Background(), indexModels)
	if err != nil {
		return nil, err
	}

	return &mongoUserRepository{
		client:     client,
		dbName:     dbName,
		collection: collection,
	}, nil
}

func (r *mongoUserRepository) GetUser(ctx context.Context, username string) (*domain.User, error) {
	collection := r.client.Database(r.dbName).Collection(r.collection)
	var user domain.User
	err := collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *mongoUserRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	collection := r.client.Database(r.dbName).Collection(r.collection)
	var user domain.User
	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *mongoUserRepository) SaveUser(ctx context.Context, user *domain.User) error {
	collection := r.client.Database(r.dbName).Collection(r.collection)
	_, err := collection.InsertOne(ctx, user)
	return err
}
