package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/nicodelara/uala-challenge/internal/common"
	timelineApp "github.com/nicodelara/uala-challenge/internal/timeline/application"
	timelineHTTP "github.com/nicodelara/uala-challenge/internal/timeline/infrastructure/http"
	timelineMongo "github.com/nicodelara/uala-challenge/internal/timeline/infrastructure/mongo"
	redisCache "github.com/nicodelara/uala-challenge/internal/timeline/infrastructure/redis"
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

	// Configurar Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
	})
	defer redisClient.Close()

	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Error connecting to Redis: %v", err)
	}

	// Inicializar repositorios
	timelineRepo, err := timelineMongo.NewMongoTimelineRepository(mongoClient, cfg.MongoDBName, "tweets")
	if err != nil {
		log.Fatalf("Error creating timeline repository: %v", err)
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

	// Inicializar repositorio de caché con TTL
	cacheRepo := redisCache.NewRedisCacheRepository(redisClient, cfg.CacheTTL)

	// Inicializar servicios
	timelineService := timelineApp.NewTimelineService(
		timelineRepo,
		userChecker,
		cacheRepo,
	)

	// Configurar router
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(logger.GinLogger())

	// Configurar handlers
	timelineHandler := timelineHTTP.NewTimelineHandler(timelineService)

	// Rutas
	timelineGroup := router.Group("/timeline")
	{
		timelineGroup.GET("/:username", timelineHandler.GetTimeline)
	}

	// Configurar servidor HTTP
	srv := &http.Server{
		Addr:           ":" + cfg.TimelinePort,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// Iniciar servidor en una goroutine
	go func() {
		logger.Info("Timeline service running on port " + cfg.TimelinePort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
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
