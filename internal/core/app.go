package core

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yakuter/ugin/internal/config"
	"github.com/yakuter/ugin/internal/domain"
	httpHandler "github.com/yakuter/ugin/internal/handler/http"
	"github.com/yakuter/ugin/internal/repository/gormrepo"
	"github.com/yakuter/ugin/internal/service"
	"github.com/yakuter/ugin/pkg/logger"
	"gorm.io/gorm"
)

// App represents the application
type App struct {
	config *config.Config
	logger *logger.Logger
	db     *gorm.DB
	server *http.Server
}

// New creates a new application instance
func New() (*App, error) {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Initialize logger
	appLogger := logger.New("info")

	appLogger.Info("initializing application")

	// Initialize database
	db, err := InitDatabase(cfg, appLogger)
	if err != nil {
		appLogger.Close()
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// Auto migrate
	if err := autoMigrate(db); err != nil {
		appLogger.Close()
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return &App{
		config: cfg,
		logger: appLogger,
		db:     db,
	}, nil
}

// Run starts the application
func (a *App) Run() error {
	defer a.logger.Close()

	// Initialize repositories
	postRepo := gormrepo.NewPostRepository(a.db)
	userRepo := gormrepo.NewUserRepository(a.db)

	// Initialize services
	authConfig := &service.AuthConfig{
		JWTSecret:            a.config.JWT.Secret,
		AccessTokenDuration:  a.config.JWT.AccessTokenDuration,
		RefreshTokenDuration: a.config.JWT.RefreshTokenDuration,
	}
	postService := service.NewPostService(postRepo, a.logger)
	authService := service.NewAuthService(userRepo, authConfig, a.logger)

	// Initialize handlers
	postHandler := httpHandler.NewPostHandler(postService)
	authHandler := httpHandler.NewAuthHandler(authService)

	// Setup router
	router := SetupRouter(a.config, postHandler, authHandler, authService, a.logger)

	// Create server
	addr := fmt.Sprintf("%s:%s", a.config.Server.Host, a.config.Server.Port)
	a.server = &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		a.logger.Info("server starting", "address", addr)
		fmt.Printf("🚀 Server is starting at %s\n", addr)
		fmt.Printf("📝 Check logs for details\n")

		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.logger.Fatal("server failed to start", "error", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	return a.waitForShutdown()
}

// waitForShutdown waits for interrupt signal and performs graceful shutdown
func (a *App) waitForShutdown() error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	a.logger.Info("shutting down server...")
	fmt.Println("\n⏳ Shutting down gracefully...")

	// Give outstanding requests 5 seconds to complete
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.server.Shutdown(ctx); err != nil {
		a.logger.Error("server forced to shutdown", "error", err)
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	a.logger.Info("server exited successfully")
	fmt.Println("✅ Server stopped")
	return nil
}

// autoMigrate runs database migrations
func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&domain.Post{},
		&domain.Tag{},
		&domain.User{},
	)
}

// setupAccessLog configures the access log file for Gin
func setupAccessLog(appLogger *logger.Logger) {
	accessLogFile, err := os.OpenFile("ugin.access.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		appLogger.Error("failed to open access log file", "error", err)
		return
	}
	gin.DefaultWriter = io.MultiWriter(accessLogFile)
}

// setGinMode sets the Gin mode based on environment variable
func setGinMode() {
	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "" {
		ginMode = gin.ReleaseMode
	}
	gin.SetMode(ginMode)
}
