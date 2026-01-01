package main

import (
	"context"
	"log"
	"os"

	"github.com/flashtix/server/db"
	"github.com/flashtix/server/internal/handlers"
	"github.com/flashtix/server/internal/middleware"
	"github.com/flashtix/server/internal/repository/postgres"
	"github.com/flashtix/server/internal/repository/redis"
	"github.com/flashtix/server/internal/services"
	"github.com/gin-gonic/gin"
	redisClient "github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Initialize Prisma Client
	client := db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer func() {
		if err := client.Prisma.Disconnect(); err != nil {
			log.Fatal("Failed to disconnect from database:", err)
		}
	}()

	// Run migrations (optional, for development)
	if _, err := client.Prisma.Raw.ExecuteRaw("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Exec(context.Background()); err != nil {
		log.Println("Warning: Could not create uuid extension:", err)
	}

	// Redis connection
	rdb := redisClient.NewClient(&redisClient.Options{
		Addr: os.Getenv("REDIS_URL"),
	})

	// Repositories
	eventRepo := postgres.NewEventRepository(client)
	ticketRepo := postgres.NewTicketRepository(client)
	seatLockRepo := redis.NewSeatLockRepository(rdb)

	// Services
	ticketService := services.NewTicketService(ticketRepo, eventRepo, seatLockRepo)

	// Handlers
	ticketHandler := handlers.NewTicketHandler(ticketService)
	eventHandler := handlers.NewEventHandler(eventRepo)

	// Router
	r := gin.Default()

	// Middleware
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.LoggingMiddleware())

	// Routes
	api := r.Group("/api")
	{
		api.GET("/events", eventHandler.GetEvents)
		api.POST("/events", eventHandler.CreateEvent)

		auth := api.Group("")
		auth.Use(middleware.AuthMiddleware(os.Getenv("JWT_SECRET")))
		{
			auth.POST("/tickets/reserve", ticketHandler.ReserveSeat)
			auth.POST("/tickets/confirm", ticketHandler.ConfirmPurchase)
		}
	}

	log.Println("Server starting on :8080")
	r.Run(":8080")
}
