package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type CreateOrderRequest struct {
	OrderID string      `json:"order_id" binding:"required"`
	Items   []OrderItem `json:"items" binding:"required,min=1"`
}

type OrderItem struct {
	ExternalID string `json:"external_id" binding:"required"`
	Name       string `json:"name" binding:"required"`
	Price      int64  `json:"price" binding:"gte=0"`
	Quantity   int    `json:"quantity" binding:"gt=0"`
}

type CreateOrderResponse struct {
	ExternalOrderID string `json:"external_order_id"`
	Status          string `json:"status"`
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	router.POST("/orders", createOrder)

	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "8081"
	}

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	logger.Info("restaurant mock started", "port", port)

	if err := server.ListenAndServe(); err != nil {
		logger.Error("server stopped", "error", err)
		os.Exit(1)
	}
}

func createOrder(c *gin.Context) {
	var request CreateOrderRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid order",
		})
		return
	}

	// Имитируем успешную передачу заказа заведению.
	c.JSON(http.StatusCreated, CreateOrderResponse{
		ExternalOrderID: "restaurant-" + request.OrderID,
		Status:          "accepted",
	})
}
