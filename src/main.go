package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

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
	channelRepo := repositories.NewChannelRepository(database.DB)
	externalUserRepo := repositories.NewExternalUserRepository(database.DB)
	conversationRepo := repositories.NewConversationRepository(database.DB)
	messageRepo := repositories.NewMessageRepository(database.DB)
	webhookEventRepo := repositories.NewWebhookEventRepository(database.DB)

	// Initialize services
	orgService := services.NewOrganizationService(orgRepo, emitter)
	channelService := services.NewChannelService(channelRepo, emitter)
	messageService := services.NewMessageService(messageRepo, conversationRepo, externalUserRepo, emitter)
	conversationService := services.NewConversationService(conversationRepo, emitter)
	webhookService := services.NewWebhookService(webhookEventRepo, messageService)

	// Initialize handlers
	orgHandler := handlers.NewOrganizationHandler(orgService)
	channelHandler := handlers.NewChannelHandler(channelService)
	messageHandler := handlers.NewMessageHandler(messageService)
	conversationHandler := handlers.NewConversationHandler(conversationService)
	webhookHandler := handlers.NewWebhookHandler(webhookService)

	// Setup Chi router
	r := chi.NewRouter()

	// Global middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(custommiddleware.SetupCORS([]string{"http://localhost:3000", "http://localhost:5173"}))

	// Public routes (no JWT required)
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	// Webhook routes (platform-specific authentication would be added per platform)
	r.Route("/api/v1/webhooks", func(r chi.Router) {
		r.Post("/{channelId}/{platform}", webhookHandler.HandleWebhook)
	})

	// API routes with JWT authentication
	r.Route("/api/v1", func(r chi.Router) {
		r.Use(custommiddleware.SetupJWT(custommiddleware.JWTConfig{
			Secret: cfg.JWT.Secret,
		}))

		// Organization routes
		r.Post("/organizations", orgHandler.Create)
		r.Get("/organizations", orgHandler.List)
		r.Get("/organizations/{id}", orgHandler.GetByID)
		r.Get("/organizations/slug/{slug}", orgHandler.GetBySlug)
		r.Patch("/organizations/{id}", orgHandler.Update)
		r.Delete("/organizations/{id}", orgHandler.Delete)

		// Channel routes
		r.Post("/channels", channelHandler.Create)
		r.Get("/channels/{id}", channelHandler.GetByID)
		r.Get("/organizations/{orgId}/channels", channelHandler.ListByOrganization)
		r.Patch("/channels/{id}", channelHandler.Update)
		r.Patch("/channels/{id}/status", channelHandler.UpdateStatus)
		r.Delete("/channels/{id}", channelHandler.Delete)

		// Conversation routes
		r.Get("/conversations/{id}", conversationHandler.GetByID)
		r.Get("/channels/{channelId}/conversations", conversationHandler.List)
		r.Post("/conversations/{id}/assign", conversationHandler.Assign)
		r.Patch("/conversations/{id}/status", conversationHandler.UpdateStatus)
		r.Patch("/conversations/{id}/priority", conversationHandler.UpdatePriority)

		// Message routes
		r.Post("/conversations/{id}/messages", messageHandler.SendMessage)
		r.Get("/conversations/{id}/messages", messageHandler.GetHistory)
		r.Post("/messages/{id}/delivered", messageHandler.MarkDelivered)
		r.Post("/messages/{id}/read", messageHandler.MarkRead)
	})

	// Start server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Server.Port),
		Handler: r,
	}

	go func() {
		log.Printf("Starting server on %s (env: %s)", server.Addr, cfg.Server.Env)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Server error: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
}
