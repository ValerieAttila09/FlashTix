package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/flashtix/server/db"
	"github.com/flashtix/server/internal/handlers"
	"github.com/flashtix/server/internal/middleware"
	"github.com/flashtix/server/internal/repository/postgres"
	"github.com/flashtix/server/internal/repository/redis"
	"github.com/flashtix/server/internal/services"
	"github.com/gin-gonic/gin"
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

	// Redis connection (Upstash)
	redisURL := os.Getenv("UPSTASH_REDIS_REST_URL")
	redisToken := os.Getenv("UPSTASH_REDIS_REST_TOKEN")

	if redisURL == "" || redisToken == "" {
		log.Println("Warning: UPSTASH_REDIS_REST_URL and UPSTASH_REDIS_REST_TOKEN not set")
		log.Println("Redis features (seat locking) will not work properly")
		redisURL = "https://dummy.upstash.io" // dummy URL to prevent panic
		redisToken = "dummy"
	}

	// Repositories
	eventRepo := postgres.NewEventRepository(client)
	ticketRepo := postgres.NewTicketRepository(client)
	seatLockRepo := redis.NewSeatLockRepository(redisURL, redisToken)

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
		// Root endpoint untuk informasi API
		api.GET("/", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "FlashTix API Server",
				"version": "1.0.0",
				"endpoints": gin.H{
					"GET /api/":                 "API information",
					"GET /api/events":           "Get all events",
					"POST /api/events":          "Create new event",
					"GET /api/redis-test":       "Test Redis connection",
					"POST /api/tickets/reserve": "Reserve a ticket seat (requires auth)",
					"POST /api/tickets/confirm": "Confirm ticket purchase (requires auth)",
				},
			})
		})

		api.GET("/events", eventHandler.GetEvents)
		api.POST("/events", eventHandler.CreateEvent)

		// Test Redis endpoint
		api.GET("/redis-test", func(c *gin.Context) {
			ctx := context.Background()
			testValue := "Hello from Upstash Redis!"

			// Test SET
			err := seatLockRepo.LockSeat(ctx, "test", "test", testValue, 60*time.Second)
			if err != nil {
				c.JSON(500, gin.H{"error": "Failed to set Redis key", "details": err.Error()})
				return
			}

			// Test GET
			retrievedValue, err := seatLockRepo.IsSeatLocked(ctx, "test", "test")
			if err != nil {
				c.JSON(500, gin.H{"error": "Failed to get Redis key", "details": err.Error()})
				return
			}

			c.JSON(200, gin.H{
				"message":         "Redis Upstash connection successful!",
				"set_value":       testValue,
				"retrieved_value": retrievedValue,
			})
		})

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
