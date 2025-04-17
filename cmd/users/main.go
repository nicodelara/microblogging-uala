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
	"github.com/nicodelara/microblogging-uala/internal/users/application"
	userhttp "github.com/nicodelara/microblogging-uala/internal/users/infrastructure/http"
	usermongo "github.com/nicodelara/microblogging-uala/internal/users/infrastructure/mongo"
	"github.com/nicodelara/microblogging-uala/pkg/config"
	"github.com/nicodelara/microblogging-uala/pkg/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Inicializar logger
	logger.Init()

	// Cargar configuración
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

	// Inicializar repositorio de usuarios
	userRepo, err := usermongo.NewMongoUserRepository(mongoClient, cfg.MongoDBName, "users")
	if err != nil {
		log.Fatalf("Error creating user repository: %v", err)
	}

	// Inicializar repositorio de follows
	followRepo, err := usermongo.NewMongoFollowRepository(mongoClient, cfg.MongoDBName, "follows")
	if err != nil {
		log.Fatalf("Error creating follow repository: %v", err)
	}

	// Crear servicio de usuarios
	userSvc := application.NewUserService(userRepo, followRepo)

	// Configurar router
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(logger.GinLogger())

	// Configurar handlers
	userHandler := userhttp.NewUserHandler(userSvc)

	// Configurar rutas
	usersGroup := router.Group("/users")
	{
		usersGroup.POST("", userHandler.CreateUser)
		usersGroup.POST("/:username/follow", userHandler.FollowUser)
	}

	// Configurar servidor HTTP
	srv := &stdhttp.Server{
		Addr:           ":" + cfg.UsersPort,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// Iniciar servidor en una goroutine
	go func() {
		logger.Info("Users service running on port " + cfg.UsersPort)
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
