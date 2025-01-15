package handlers

import (
	"net/http"
	"strconv"

	"github.com/simon4545/goshop/database"
	"github.com/simon4545/goshop/models"

	"github.com/labstack/echo/v4"
)

type OrderRequest struct {
	Products   []models.Product `json:"products"`
	Status     string           `json:"status"`
	FirstName  string           `json:"first_name"`
	LastName   string           `json:"last_name"`
	Phone      string           `json:"phone"`
	Email      string           `json:"email"`
	Address    string           `json:"address"`
	PostalCode string           `json:"postal_code"`
}

func CreateOrder(c echo.Context) error {
	var orderRequest OrderRequest
	if err := c.Bind(&orderRequest); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	order := models.Order{
		Status:     orderRequest.Status,
		FirstName:  orderRequest.FirstName,
		LastName:   orderRequest.LastName,
		Phone:      orderRequest.Phone,
		Email:      orderRequest.Email,
		Address:    orderRequest.Address,
		PostalCode: orderRequest.PostalCode,
	}

	if err := database.DB.Create(&order).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create order"})
	}

	// Add products to order
	for _, product := range orderRequest.Products {
		orderProduct := models.Item{
			OrderID:   order.ID,
			ProductID: product.ID,
		}
		database.DB.Create(&orderProduct)
	}

	return c.JSON(http.StatusCreated, order)
}

func GetOrder(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	var order models.Order
	database.DB.Preload("Products").First(&order, id)
	return c.JSON(http.StatusOK, order)
}

func UpdateOrderStatus(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	var order models.Order
	database.DB.First(&order, id)
	if order.ID == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Order not found"})
	}

	var statusUpdate struct {
		Status string `json:"status"`
	}
	if err := c.Bind(&statusUpdate); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	order.Status = statusUpdate.Status
	database.DB.Save(&order)
	return c.JSON(http.StatusOK, order)
}

func DeleteOrder(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	var order models.Order
	database.DB.First(&order, id)
	if order.ID == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Order not found"})
	}
	database.DB.Delete(&order)
	return c.JSON(http.StatusOK, map[string]string{"message": "Order deleted"})
}
