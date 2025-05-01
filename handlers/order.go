package handlers

import (
	"net/http"
	"strconv"

	"github.com/simon4545/goshop/database"
	"github.com/simon4545/goshop/models"

	"github.com/labstack/echo/v4"
)

type OrderRequest struct {
	Products   []models.Product   `json:"products"`
	CustomerID uint               `json:"customer_id"`
	AgentID    uint               `json:"agent_id"` // 订单所属代理商
	OrderNo    string             `gorm:"size:50;unique" json:"order_no"`
	TotalPrice float64            `gorm:"type:decimal(10,2);not null" json:"total_price"`
	Status     models.OrderStatus `gorm:"size:20" json:"status"`
	Address    string             `gorm:"size:255" json:"address"`
	Phone      string             `gorm:"size:20" json:"phone"`
	Remark     string             `gorm:"size:255" json:"remark"`
}

func CreateOrder(c echo.Context) error {
	var orderRequest OrderRequest
	if err := c.Bind(&orderRequest); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	order := models.Order{
		Status:     models.OrderStatus(orderRequest.Status),
		CustomerID: 1,
		AgentID:    1,
		OrderNo:    orderRequest.Phone,
		TotalPrice: 1,
		Address:    orderRequest.Address,
		Remark:     orderRequest.Remark,
	}

	if err := database.DB.Create(&order).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create order"})
	}

	// Add products to order
	for _, product := range orderRequest.Products {
		orderProduct := models.OrderItem{
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

	order.Status = models.OrderStatus(statusUpdate.Status)
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
