package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	TweetsPort   string
	UsersPort    string
	TimelinePort string
	MongoURI     string
	MongoDBName  string
	UsersDBName  string
	RedisAddr    string
	KafkaBrokers string
	CacheTTL     time.Duration
}

func LoadConfig() (*Config, error) {
	tweetsPort := os.Getenv("TWEETS_PORT")
	if tweetsPort == "" {
		tweetsPort = "8081"
	}

	usersPort := os.Getenv("USERS_PORT")
	if usersPort == "" {
		usersPort = "8082"
	}

	timelinePort := os.Getenv("TIMELINE_PORT")
	if timelinePort == "" {
		timelinePort = "8083"
	}

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	mongoDBName := os.Getenv("MONGO_DB_NAME")
	if mongoDBName == "" {
		mongoDBName = "twitter"
	}

	usersDBName := os.Getenv("USERS_DB_NAME")
	if usersDBName == "" {
		usersDBName = "twitter"
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	kafkaBrokers := os.Getenv("KAFKA_BROKERS")
	if kafkaBrokers == "" {
		kafkaBrokers = "localhost:9092"
	}

	cacheTTLStr := os.Getenv("CACHE_TTL_SECONDS")
	cacheTTL := 30 * time.Second // valor por defecto
	if cacheTTLStr != "" {
		if seconds, err := strconv.Atoi(cacheTTLStr); err == nil {
			cacheTTL = time.Duration(seconds) * time.Second
		}
	}

	return &Config{
		TweetsPort:   tweetsPort,
		UsersPort:    usersPort,
		TimelinePort: timelinePort,
		MongoURI:     mongoURI,
		MongoDBName:  mongoDBName,
		UsersDBName:  usersDBName,
		RedisAddr:    redisAddr,
		KafkaBrokers: kafkaBrokers,
		CacheTTL:     cacheTTL,
	}, nil
}
