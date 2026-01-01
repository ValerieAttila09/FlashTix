package main

import (
	"log"
	"os"

	"github.com/flashtix/server/internal/domain"
	"github.com/flashtix/server/internal/handlers"
	"github.com/flashtix/server/internal/middleware"
	"github.com/flashtix/server/internal/repository/postgres"
	"github.com/flashtix/server/internal/repository/redis"
	"github.com/flashtix/server/internal/services"
	"github.com/gin-gonic/gin"
	redisClient "github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	postgresDriver "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Database connection
	db, err := gorm.Open(postgresDriver.Open(os.Getenv("DATABASE_URL")), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migrate
	db.AutoMigrate(&domain.Event{}, &domain.Ticket{}, &domain.User{})

	// Redis connection
	rdb := redisClient.NewClient(&redisClient.Options{
		Addr: os.Getenv("REDIS_URL"),
	})

	// Repositories
	eventRepo := postgres.NewEventRepository(db)
	ticketRepo := postgres.NewTicketRepository(db)
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
