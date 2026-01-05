package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github/sarthak-pokharel/sqlite-d1-gochat/src/config"
	"github/sarthak-pokharel/sqlite-d1-gochat/src/database"
	"github/sarthak-pokharel/sqlite-d1-gochat/src/events"
	"github/sarthak-pokharel/sqlite-d1-gochat/src/handlers"
	custommiddleware "github/sarthak-pokharel/sqlite-d1-gochat/src/middleware"
	"github/sarthak-pokharel/sqlite-d1-gochat/src/repositories"
	"github/sarthak-pokharel/sqlite-d1-gochat/src/services"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	if err := database.InitDB(cfg.Database.Path); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	// Run GORM auto-migrations
	if err := database.AutoMigrate(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize Redis event emitter
	emitter, err := events.NewRedisEmitter(
		fmt.Sprintf("redis://%s:%d/%d", cfg.Redis.Host, cfg.Redis.Port, cfg.Redis.DB),
		cfg.Redis.Enabled,
	)
	if err != nil {
		log.Fatalf("Failed to initialize event emitter: %v", err)
	}
	defer emitter.Close()

	// Initialize repositories
	orgRepo := repositories.NewOrganizationRepository(database.DB)
	_ = repositories.NewChannelRepository(database.DB) // TODO: wire up channel service
	externalUserRepo := repositories.NewExternalUserRepository(database.DB)
	conversationRepo := repositories.NewConversationRepository(database.DB)
	messageRepo := repositories.NewMessageRepository(database.DB)
	webhookEventRepo := repositories.NewWebhookEventRepository(database.DB)

	// Initialize services
	orgService := services.NewOrganizationService(orgRepo, emitter)
	messageService := services.NewMessageService(messageRepo, conversationRepo, externalUserRepo, emitter)
	conversationService := services.NewConversationService(conversationRepo, emitter)
	webhookService := services.NewWebhookService(webhookEventRepo, messageService)

	// Initialize handlers
	orgHandler := handlers.NewOrganizationHandler(orgService)
	messageHandler := handlers.NewMessageHandler(messageService)
	conversationHandler := handlers.NewConversationHandler(conversationService)
	webhookHandler := handlers.NewWebhookHandler(webhookService)

	// Setup Echo server
	e := echo.New()
	e.HideBanner = true
	e.HTTPErrorHandler = custommiddleware.CustomErrorHandler

	// Global middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(custommiddleware.SetupCORS([]string{"http://localhost:3000", "http://localhost:5173"}))

	// Public routes (no JWT required)
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "ok"})
	})

	// Webhook routes (platform-specific authentication would be added per platform)
	webhooks := e.Group("/api/v1/webhooks")
	webhooks.POST("/:channelId/:platform", webhookHandler.HandleWebhook)

	// API routes with JWT authentication
	api := e.Group("/api/v1")
	api.Use(custommiddleware.SetupJWT(custommiddleware.JWTConfig{
		Secret: cfg.JWT.Secret,
		SkipRoutes: []string{
			"/api/v1/webhooks",
		},
	}))

	// Organization routes
	api.POST("/organizations", orgHandler.Create)
	api.GET("/organizations", orgHandler.List)
	api.GET("/organizations/:id", orgHandler.GetByID)
	api.GET("/organizations/slug/:slug", orgHandler.GetBySlug)
	api.PATCH("/organizations/:id", orgHandler.Update)
	api.DELETE("/organizations/:id", orgHandler.Delete)

	// Conversation routes
	api.GET("/conversations/:id", conversationHandler.GetByID)
	api.GET("/channels/:channelId/conversations", conversationHandler.List)
	api.POST("/conversations/:id/assign", conversationHandler.Assign)
	api.PATCH("/conversations/:id/status", conversationHandler.UpdateStatus)
	api.PATCH("/conversations/:id/priority", conversationHandler.UpdatePriority)

	// Message routes
	api.POST("/conversations/:id/messages", messageHandler.SendMessage)
	api.GET("/conversations/:id/messages", messageHandler.GetHistory)
	api.POST("/messages/:id/delivered", messageHandler.MarkDelivered)
	api.POST("/messages/:id/read", messageHandler.MarkRead)

	// Start server
	go func() {
		addr := fmt.Sprintf(":%s", cfg.Server.Port)
		log.Printf("Starting server on %s (env: %s)", addr, cfg.Server.Env)
		if err := e.Start(addr); err != nil {
			log.Printf("Server error: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
}
