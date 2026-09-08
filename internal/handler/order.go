package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-goalexey-gsv-15625069/internal/model"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-goalexey-gsv-15625069/internal/repository"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-goalexey-gsv-15625069/internal/service"
)

type OrderHandler struct {
	service *service.OrderService
}

func NewOrderHandler(service *service.OrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

type CreateOrderRequest struct {
	RestaurantID    string            `json:"restaurant_id" binding:"required"`
	DeliveryAddress string            `json:"delivery_address" binding:"required"`
	Items           []CreateOrderItem `json:"items" binding:"required,min=1"`
}

type CreateOrderItem struct {
	MenuItemID string `json:"menu_item_id" binding:"required"`
	Quantity   int    `json:"quantity" binding:"required,gt=0"`
}

// Create godoc
// @Summary Create order
// @Description Create a new food delivery order. The order is sent to the restaurant service for acceptance.
// @Tags orders
// @Accept json
// @Produce json
// @Param X-User-ID header string true "User UUID" example(99999999-9999-9999-9999-999999999999)
// @Param request body CreateOrderRequest true "Order data"
// @Success 201 {object} model.Order
// @Failure 400 {object} map[string]string
// @Failure 502 {object} map[string]string
// @Router /orders [post]
func (h *OrderHandler) Create(c *gin.Context) {
	var request CreateOrderRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	restaurantID, err := uuid.Parse(request.RestaurantID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid restaurant_id",
		})
		return
	}

	userID, err := uuid.Parse(c.GetHeader("X-User-ID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "X-User-ID must be a valid UUID",
		})
		return
	}

	items := make([]repository.CreateOrderItem, 0, len(request.Items))

	for _, item := range request.Items {
		menuItemID, err := uuid.Parse(item.MenuItemID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid menu_item_id",
			})
			return
		}

		items = append(items, repository.CreateOrderItem{
			MenuItemID: menuItemID,
			Quantity:   item.Quantity,
		})
	}

	order, err := h.service.Create(
		c.Request.Context(),
		service.CreateOrderInput{
			RestaurantID:    restaurantID,
			UserID:          userID,
			DeliveryAddress: request.DeliveryAddress,
			Items:           items,
		},
	)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, order)
}

// GetByID godoc
// @Summary Get order
// @Description Get an order by ID including its items.
// @Tags orders
// @Produce json
// @Param id path string true "Order UUID" format(uuid)
// @Success 200 {object} model.Order
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /orders/{id} [get]
func (h *OrderHandler) GetByID(c *gin.Context) {
	orderID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid order id",
		})
		return
	}

	order, err := h.service.GetByID(
		c.Request.Context(),
		orderID,
	)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "order not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get order",
		})
		return
	}

	c.JSON(http.StatusOK, order)
}

// GetByUserID godoc
// @Summary Get user orders
// @Description Get all orders created by a user.
// @Tags orders
// @Produce json
// @Param userId path string true "User UUID" format(uuid)
// @Success 200 {array} model.Order
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/{userId}/orders [get]
func (h *OrderHandler) GetByUserID(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid user id",
		})
		return
	}

	orders, err := h.service.GetByUserID(
		c.Request.Context(),
		userID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get user orders",
		})
		return
	}

	c.JSON(http.StatusOK, orders)
}

// Make sure imported model package is used by swag annotations.
var _ model.Order
