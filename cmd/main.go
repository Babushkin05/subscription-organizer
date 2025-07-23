package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/Babushkin05/subscription-organizer/docs"
	"github.com/Babushkin05/subscription-organizer/internal/application/usecase"
	"github.com/Babushkin05/subscription-organizer/internal/config"
	httpService "github.com/Babushkin05/subscription-organizer/internal/infrastructure/delivery/http"
	"github.com/Babushkin05/subscription-organizer/internal/infrastructure/repository/postgres"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	// Load config
	cfg := config.MustLoad()
	log.Println("Config loaded successfully")

	// Connect to DB
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.DataBase.Host,
		cfg.DataBase.Port,
		cfg.DataBase.User,
		cfg.DataBase.Password,
		cfg.DataBase.Name)
	log.Println("DSN:", dsn)
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}
	defer db.Close()
	log.Println("Connected to DB")

	// Init repository
	subRepo := postgres.NewSubscriptionRepository(db)

	// Init service
	subService := usecase.NewSubscriptionService(subRepo)

	// Init Gin router
	r := gin.Default()

	// Init handler
	subHandler := httpService.NewSubscriptionHandler(subService)

	// Register routes
	httpService.RegisterRoutes(r, subHandler)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Run server
	addr := ":" + strconv.Itoa(cfg.Server.Port)
	log.Printf("Starting server at %s\n", addr)
	if err := r.Run(addr); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}
