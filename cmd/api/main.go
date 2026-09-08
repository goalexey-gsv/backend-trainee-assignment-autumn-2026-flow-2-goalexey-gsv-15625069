package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pressly/goose/v3"

	httpSwagger "github.com/swaggo/http-swagger/v2"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-goalexey-gsv-15625069/internal/handler"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-goalexey-gsv-15625069/internal/integration/restaurant"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-goalexey-gsv-15625069/internal/repository/postgres"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-goalexey-gsv-15625069/internal/service"
)

// @title Kitchen API
// @version 1.0
// @description MVP API for restaurants and food delivery orders.
// @host localhost:8080
// @BasePath /api/v1
func main() {
	logger := slog.New(
		slog.NewTextHandler(os.Stdout, nil),
	)
	slog.SetDefault(logger)

	databaseURL := getEnv(
		"DATABASE_URL",
		"postgres://kitchen:kitchen@localhost:5432/kitchen?sslmode=disable",
	)

	httpPort := getEnv("HTTP_PORT", "8080")

	restaurantServiceURL := getEnv(
		"RESTAURANT_SERVICE_URL",
		"http://localhost:8081",
	)

	ctx := context.Background()

	db, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		logger.Error(
			"failed to create database pool",
			"error",
			err,
		)
		os.Exit(1)
	}
	defer db.Close()

	if err := waitForDatabase(ctx, db); err != nil {
		logger.Error(
			"database is unavailable",
			"error",
			err,
		)
		os.Exit(1)
	}

	if err := runMigrations(db); err != nil {
		logger.Error(
			"failed to run migrations",
			"error",
			err,
		)
		os.Exit(1)
	}

	// Repositories.

	restaurantRepository := postgres.NewRestaurantRepository(db)
	orderRepository := postgres.NewOrderRepository(db)

	// Integrations.

	restaurantClient := restaurant.NewClient(
		restaurantServiceURL,
	)

	// Services.

	restaurantService := service.NewRestaurantService(
		restaurantRepository,
	)

	orderService := service.NewOrderService(
		orderRepository,
		restaurantClient,
	)

	// Handlers.

	restaurantHandler := handler.NewRestaurantHandler(
		restaurantService,
	)

	orderHandler := handler.NewOrderHandler(
		orderService,
	)

	// Router.

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		if err := db.Ping(c.Request.Context()); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "unhealthy",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	api := router.Group("/api/v1")
	{
		// Restaurants.

		api.GET(
			"/restaurants",
			restaurantHandler.List,
		)

		api.GET(
			"/restaurants/:id",
			restaurantHandler.GetByID,
		)

		api.GET(
			"/restaurants/:id/menu",
			restaurantHandler.GetMenu,
		)

		// Orders.

		api.POST(
			"/orders",
			orderHandler.Create,
		)

		api.GET(
			"/orders/:id",
			orderHandler.GetByID,
		)

		api.GET(
			"/users/:userId/orders",
			orderHandler.GetByUserID,
		)
	}

	router.GET(
		"/swagger/*any",
		gin.WrapH(httpSwagger.Handler(
			httpSwagger.URL("/swagger/doc.json"),
		)),
	)

	server := &http.Server{
		Addr:              ":" + httpPort,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		logger.Info(
			"HTTP server started",
			"port",
			httpPort,
		)

		if err := server.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			logger.Error(
				"HTTP server stopped",
				"error",
				err,
			)
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)

	signal.Notify(
		stop,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	<-stop

	logger.Info("shutting down HTTP server")

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error(
			"failed to shutdown HTTP server",
			"error",
			err,
		)
		os.Exit(1)
	}

	logger.Info("HTTP server stopped")
}

func runMigrations(db *pgxpool.Pool) error {
	sqlDB, err := goose.OpenDBWithDriver(
		"pgx",
		db.Config().ConnString(),
	)
	if err != nil {
		return err
	}
	defer sqlDB.Close()

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	return goose.Up(
		sqlDB,
		"/app/migrations",
	)
}

func waitForDatabase(
	ctx context.Context,
	db *pgxpool.Pool,
) error {
	const (
		maxAttempts = 30
		retryDelay  = time.Second
	)

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		if err := db.Ping(ctx); err == nil {
			return nil
		}

		time.Sleep(retryDelay)
	}

	return errors.New(
		"database connection timeout",
	)
}

func getEnv(
	key string,
	fallback string,
) string {
	value := os.Getenv(key)

	if value == "" {
		return fallback
	}

	return value
}
