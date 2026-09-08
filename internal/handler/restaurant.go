package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-goalexey-gsv-15625069/internal/model"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-goalexey-gsv-15625069/internal/service"
)

type RestaurantHandler struct {
	service *service.RestaurantService
}

func NewRestaurantHandler(service *service.RestaurantService) *RestaurantHandler {
	return &RestaurantHandler{service: service}
}

// List godoc
// @Summary List restaurants
// @Description Get a list of all active restaurants.
// @Tags restaurants
// @Produce json
// @Success 200 {array} model.Restaurant
// @Failure 500 {object} map[string]string
// @Router /restaurants [get]
func (h *RestaurantHandler) List(c *gin.Context) {
	restaurants, err := h.service.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get restaurants",
		})
		return
	}

	c.JSON(http.StatusOK, restaurants)
}

// GetByID godoc
// @Summary Get restaurant
// @Description Get an active restaurant by ID.
// @Tags restaurants
// @Produce json
// @Param id path string true "Restaurant UUID" format(uuid)
// @Success 200 {object} model.Restaurant
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /restaurants/{id} [get]
func (h *RestaurantHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid restaurant id",
		})
		return
	}

	restaurant, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "restaurant not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get restaurant",
		})
		return
	}

	c.JSON(http.StatusOK, restaurant)
}

// GetMenu godoc
// @Summary Get restaurant menu
// @Description Get all available menu items for a restaurant.
// @Tags restaurants
// @Produce json
// @Param id path string true "Restaurant UUID" format(uuid)
// @Success 200 {array} model.MenuItem
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /restaurants/{id}/menu [get]
func (h *RestaurantHandler) GetMenu(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid restaurant id",
		})
		return
	}

	items, err := h.service.GetMenu(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get menu",
		})
		return
	}

	c.JSON(http.StatusOK, items)
}

// Make sure imported model package is used by swag annotations.
var _ model.Restaurant
