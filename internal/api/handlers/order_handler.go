package handlers

import (
	"applicationDesignTest/internal/api/dto"
	"applicationDesignTest/internal/model"
	"applicationDesignTest/internal/services"
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type OrderHandler struct {
	bookingManager services.BookingManager
}

func NewOrderHandler(bookingManager services.BookingManager) *OrderHandler {
	return &OrderHandler{
		bookingManager: bookingManager,
	}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var request dto.OrderRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if request.CheckIn.After(request.CheckOut) || request.CheckIn.Equal(request.CheckOut) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "check_in date must be before check_out date"})
		return
	}

	order := request.ToModel()
	order.OrderId = uuid.New().String()

	orders, err := h.bookingManager.HandleOrders(context.Background(), []model.Order{order})
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, model.ErrInsufficientQuota) {
			status = http.StatusConflict
		}
		if errors.Is(err, model.ErrNotFound) {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	responses := dto.FromModelList(orders)
	c.JSON(http.StatusCreated, responses)
}

func (h *OrderHandler) RegisterRoutes(router *gin.Engine) {
	router.POST("/api/orders", h.CreateOrder)
}
