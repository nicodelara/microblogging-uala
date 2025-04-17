package main

import (
	"context"
	"log"
	stdhttp "net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nicodelara/uala-challenge/internal/common"
	tweetApp "github.com/nicodelara/uala-challenge/internal/tweets/application"
	tweetHTTP "github.com/nicodelara/uala-challenge/internal/tweets/infrastructure/http"
	tweetMongo "github.com/nicodelara/uala-challenge/internal/tweets/infrastructure/mongo"
	userMongo "github.com/nicodelara/uala-challenge/internal/users/infrastructure/mongo"
	"github.com/nicodelara/uala-challenge/pkg/config"
	"github.com/nicodelara/uala-challenge/pkg/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	logger.Init()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Configurar MongoDB
	mongoClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
	}
	defer mongoClient.Disconnect(context.Background())

	// Verificar conexión a MongoDB
	if err := mongoClient.Ping(context.Background(), nil); err != nil {
		log.Fatalf("Error verifying MongoDB connection: %v", err)
	}

	// Inicializar repositorios
	tweetRepo, err := tweetMongo.NewMongoTweetRepository(mongoClient, cfg.MongoDBName, "tweets")
	if err != nil {
		log.Fatalf("Error creating tweet repository: %v", err)
	}

	userRepo, err := userMongo.NewMongoUserRepository(mongoClient, cfg.UsersDBName, "users")
	if err != nil {
		log.Fatalf("Error creating user repository: %v", err)
	}

	followRepo, err := userMongo.NewMongoFollowRepository(mongoClient, cfg.UsersDBName, "follows")
	if err != nil {
		log.Fatalf("Error creating follow repository: %v", err)
	}

	// Crear adaptador para UserChecker
	userChecker := common.NewUserCheckerAdapter(userRepo, followRepo)

	// Inicializar servicio
	tweetService := tweetApp.NewTweetService(tweetRepo, userChecker)

	// Configurar router
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(logger.GinLogger())

	// Configurar handlers
	tweetHandler := tweetHTTP.NewTweetHandler(tweetService)

	// Rutas
	tweetsGroup := router.Group("/tweets")
	{
		tweetsGroup.POST("", tweetHandler.CreateTweet)
	}

	// Configurar servidor HTTP
	srv := &stdhttp.Server{
		Addr:           ":" + cfg.TweetsPort,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// Iniciar servidor en una goroutine
	go func() {
		logger.Info("Tweets service running on port " + cfg.TweetsPort)
		if err := srv.ListenAndServe(); err != nil && err != stdhttp.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Esperar señal de interrupción
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	// Dar tiempo para que las conexiones se cierren
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	logger.Info("Server exiting")
}
