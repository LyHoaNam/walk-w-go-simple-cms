package main

import (
	"fmt"
	"log"

	"simple-template/internal/config"
	"simple-template/internal/database"
	"simple-template/internal/handler"
	"simple-template/internal/middleware"
	"simple-template/internal/repository"
	"simple-template/internal/usecase"
	"simple-template/pkg/response"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to database
	db, err := database.Connect(database.Config{
		Host:            cfg.Database.Host,
		Port:            cfg.Database.Port,
		User:            cfg.Database.User,
		Password:        cfg.Database.Password,
		Name:            cfg.Database.Name,
		MaxOpenConns:    cfg.Database.MaxOpenConns,
		MaxIdleConns:    cfg.Database.MaxIdleConns,
		ConnMaxLifetime: cfg.Database.ConnMaxLifetime,
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("✅ Connected to database successfully")

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	productRepo := repository.NewProductRepository(db)

	// Initialize usecases
	userUsecase := usecase.NewUserUsecase(userRepo)
	productUsecase := usecase.NewProductUsecase(productRepo)

	// Initialize handlers
	userHandler := handler.NewUserHandler(userUsecase)
	productHandler := handler.NewProductHandler(productUsecase)

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName: "Simple Golang API",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return response.InternalServerError(c, "Internal server error", err)
		},
	})

	// Middleware
	app.Use(recover.New())             // Recover from panics
	app.Use(cors.New())                // Enable CORS
	app.Use(middleware.Logger())       // Custom logger
	app.Use(middleware.ErrorHandler()) // Custom error handler

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return response.Success(c, fiber.Map{
			"status":   "ok",
			"database": "connected",
		}, "Service is healthy")
	})

	// API routes
	api := app.Group("/api/v1")

	// User routes
	users := api.Group("/users")
	users.Post("/", userHandler.CreateUser)      // Tạo user mới
	users.Get("/", userHandler.GetAllUsers)      // Lấy danh sách users
	users.Get("/:id", userHandler.GetUser)       // Lấy user theo ID
	users.Put("/:id", userHandler.UpdateUser)    // Cập nhật user
	users.Delete("/:id", userHandler.DeleteUser) // Xóa user

	// Product routes
	product := api.Group("/products")
	product.Post("/", productHandler.CreateProduct)

	// Start server
	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	log.Printf("🚀 Server starting on %s", addr)
	if err := app.Listen(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
