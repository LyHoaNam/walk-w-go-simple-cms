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
	customerRepo := repository.NewCustomerRepository(db)
	platformRepo := repository.NewPlatformRepository(db)
	retailStoreRepo := repository.NewRetailStoreRepository(db)
	paymentMethodsRepo := repository.NewPaymentMethodsRepository(db)
	ordersRepo := repository.NewOrdersRepository(db)

	// Initialize usecases
	userUsecase := usecase.NewUserUsecase(userRepo)
	productUsecase := usecase.NewProductUsecase(productRepo)
	customerUsecase := usecase.NewCustomerUsecase(customerRepo)
	platformUsecase := usecase.NewPlatformUsecase(platformRepo)
	retailStoreUsecase := usecase.NewRetailStoreUsecase(retailStoreRepo)
	paymentMethodsUsecase := usecase.NewPaymentMethodsUsecase(paymentMethodsRepo)
	ordersUsecase := usecase.NewOrderUseCase(ordersRepo)

	// Initialize handlers
	userHandler := handler.NewUserHandler(userUsecase)
	productHandler := handler.NewProductHandler(productUsecase)
	customerHandler := handler.NewCustomerHandler(customerUsecase)
	platformHandler := handler.NewPlatformHandler(platformUsecase)
	retailStoreHandler := handler.NewRetailStoreHandler(retailStoreUsecase)
	paymentMethodsHandler := handler.NewPaymentMethodsHandler(paymentMethodsUsecase)
	ordersHandler := handler.NewOrderHandler(ordersUsecase)
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
	users.Post("/", userHandler.CreateUser)      // Create new user
	users.Get("/", userHandler.GetAllUsers)      // Get list of users
	users.Get("/:id", userHandler.GetUser)       // Get user by ID
	users.Put("/:id", userHandler.UpdateUser)    // Update user
	users.Delete("/:id", userHandler.DeleteUser) // Delete user

	product := api.Group("/products")
	product.Get("/", productHandler.GetAll)

	product.Get("/:id", productHandler.GetByID)
	product.Post("/", productHandler.CreateProduct)
	product.Delete("/:id", productHandler.DeleteProduct)

	// Customer
	customer := api.Group("/customer")
	customer.Post("/", customerHandler.Create)
	customer.Get("/", customerHandler.GetAllCustomers)
	customer.Get("/:id", customerHandler.GetCustomer)
	customer.Put("/:id", customerHandler.UpdateCustomers)
	customer.Delete("/:id", customerHandler.DeleteCustomer)

	// platform
	platform := api.Group("/platform")
	platform.Get("/", platformHandler.GetAll)
	// retail store
	retailStore := api.Group("/retail-store")
	retailStore.Get("/", retailStoreHandler.GetAll)
	// payment methods
	paymentMethods := api.Group("/payment-methods")
	paymentMethods.Get("/", paymentMethodsHandler.GetAll)
	// orders
	orders := api.Group("/orders")
	orders.Get("/", ordersHandler.GetAll)
	orders.Post("/", ordersHandler.Create)
	orders.Put("/:id", ordersHandler.UpdateStatus)

	// Start server
	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	log.Printf("🚀 Server starting on %s", addr)
	if err := app.Listen(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
