package handlers

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/simon4545/goshop/database"
	"github.com/simon4545/goshop/models"
)

type LocationRequest struct {
	OpenID   string `json:"openid"`
	Province string `json:"province"`
	City     string `json:"city"`
}

type PhoneRequest struct {
	OpenID        string `json:"openid"`
	EncryptedData string `json:"encryptedData"`
	IV            string `json:"iv"`
	SessionKey    string `json:"session_key"`
}

type PhoneInfo struct {
	PhoneNumber     string `json:"phoneNumber"`
	PurePhoneNumber string `json:"purePhoneNumber"`
	CountryCode     string `json:"countryCode"`
}

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

func HandleLocation(c echo.Context) error {
	var req LocationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	// 这里可以添加验证逻辑，比如验证openid是否有效

	// 保存到数据库或进行其他处理
	fmt.Printf("Received location: %+v\n", req)

	// 返回成功响应
	return c.JSON(http.StatusCreated, map[string]interface{}{
		"code":    0,
		"message": "位置信息保存成功",
	})
}

func HandlePhone(c echo.Context) error {
	var req PhoneRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	// 解密手机号
	phoneInfo, err := decryptPhoneNumber(req.EncryptedData, req.SessionKey, req.IV)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	// 返回手机号信息
	return c.JSON(http.StatusCreated, map[string]interface{}{
		"code":        0,
		"message":     "手机号获取成功",
		"phoneNumber": phoneInfo.PhoneNumber,
	})
}

// 微信登录处理
func HandleWxLogin(c echo.Context) error {
	type LoginRequest struct {
		Code string `json:"code"`
	}

	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	// 调用微信接口获取session_key和openid
	resp, err := http.Get(fmt.Sprintf(
		"https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code",
		"你的小程序appid",
		"你的小程序secret",
		req.Code,
	))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}
	defer resp.Body.Close()

	var result struct {
		OpenID     string `json:"openid"`
		SessionKey string `json:"session_key"`
		ErrCode    int    `json:"errcode"`
		ErrMsg     string `json:"errmsg"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	if result.ErrCode != 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	// 生成自己的session token并返回给客户端
	token, err := GenerateToken(result.OpenID) // 假设用户 ID 是 123
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Failed to generate token",
		})
	}
	// 这里简化处理，实际应该生成并存储token
	return c.JSON(http.StatusCreated, map[string]interface{}{
		"code":        0,
		"session_key": token,
		"openid":      result.OpenID,
	})
}
func decryptPhoneNumber(encryptedData, sessionKey, iv string) (*PhoneInfo, error) {
	aesKey, err := base64.StdEncoding.DecodeString(sessionKey)
	if err != nil {
		return nil, fmt.Errorf("base64解码sessionKey失败: %v", err)
	}

	cipherText, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return nil, fmt.Errorf("base64解码encryptedData失败: %v", err)
	}

	ivBytes, err := base64.StdEncoding.DecodeString(iv)
	if err != nil {
		return nil, fmt.Errorf("base64解码iv失败: %v", err)
	}

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, fmt.Errorf("创建AES cipher失败: %v", err)
	}

	mode := cipher.NewCBCDecrypter(block, ivBytes)
	mode.CryptBlocks(cipherText, cipherText)

	// 去除填充
	cipherText = pkcs7Unpad(cipherText)

	var phoneInfo PhoneInfo
	err = json.Unmarshal(cipherText, &phoneInfo)
	if err != nil {
		return nil, fmt.Errorf("解析手机号信息失败: %v", err)
	}

	return &phoneInfo, nil
}

func pkcs7Unpad(data []byte) []byte {
	length := len(data)
	unpadding := int(data[length-1])
	return data[:(length - unpadding)]
}
