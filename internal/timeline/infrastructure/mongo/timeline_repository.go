package mongo

import (
	"context"
	"time"

	"github.com/nicodelara/microblogging-uala/internal/timeline/domain"
	"github.com/nicodelara/microblogging-uala/internal/timeline/domain/ports"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// mongoTweet es la estructura que representa un tweet en MongoDB
type mongoTweet struct {
	ID        string    `bson:"_id"`
	Username  string    `bson:"username"`
	Content   string    `bson:"content"`
	CreatedAt time.Time `bson:"createdAt"`
}

// mongoTimelineRepository es la implementación de TimelineRepository usando MongoDB.
type mongoTimelineRepository struct {
	collection *mongo.Collection
}

func NewMongoTimelineRepository(client *mongo.Client, dbName, collName string) (ports.TimelineRepository, error) {
	db := client.Database(dbName)
	collection := db.Collection(collName)
	return &mongoTimelineRepository{
		collection: collection,
	}, nil
}

// GetTweetsForUsers obtiene los tweets de una lista de usuarios
func (r *mongoTimelineRepository) GetTweetsForUsers(ctx context.Context, usernames []string, offset, limit int) ([]domain.Tweet, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	filter := bson.M{"username": bson.M{"$in": usernames}}
	opts := options.Find().
		SetSort(bson.D{{Key: "createdAt", Value: -1}}).
		SetSkip(int64(offset)).
		SetLimit(int64(limit))

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tweets []domain.Tweet
	for cursor.Next(ctx) {
		var mongoTweet mongoTweet
		if err := cursor.Decode(&mongoTweet); err != nil {
			return nil, err
		}
		tweets = append(tweets, domain.Tweet{
			ID:        mongoTweet.ID,
			Username:  mongoTweet.Username,
			Content:   mongoTweet.Content,
			CreatedAt: mongoTweet.CreatedAt,
		})
	}

	return tweets, nil
}
