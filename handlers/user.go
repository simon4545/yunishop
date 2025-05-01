package handlers

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/simon4545/goshop/database"
	"github.com/simon4545/goshop/models"
)

// 注册客户
func RegisterCustomer(c echo.Context) error {
	var input struct {
		Name       string `json:"name" binding:"required"`
		BirthDate  string `json:"birth_date" binding:"required"`
		Province   string `json:"province"`
		City       string `json:"city"`
		Address    string `json:"address"`
		Gender     string `json:"gender"`
		ReferrerID uint   `json:"referrer_id"`
		AgentID    uint   `json:"agent_id" binding:"required"`
		Phone      string `json:"phone" binding:"required"`
		Password   string `json:"password" binding:"required"`
	}

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	birthDate, err := time.Parse("2006-01-02", input.BirthDate)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	// 检查代理商是否存在
	var agent models.Agent
	if err := database.DB.First(&agent, input.AgentID).Error; err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	// 检查介绍人是否存在（如果提供了介绍人ID）
	if input.ReferrerID != 0 {
		var referrer models.Customer
		if err := database.DB.First(&referrer, input.ReferrerID).Error; err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
		}
	}

	customer := models.Customer{
		BaseUser: models.BaseUser{
			Name:       input.Name,
			BirthDate:  birthDate,
			Province:   input.Province,
			City:       input.City,
			Address:    input.Address,
			Gender:     input.Gender,
			ReferrerID: input.ReferrerID,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		AgentID: input.AgentID,
	}

	// 这里应该添加密码哈希处理
	// hashedPassword, err := utils.HashPassword(input.Password)
	// if err != nil {
	//     c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not hash password"})
	//     return
	// }
	// customer.Password = hashedPassword

	if err := database.DB.Create(&customer).Error; err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	return c.JSON(http.StatusCreated, customer)
}
