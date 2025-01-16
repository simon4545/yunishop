package handlers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo/v4"
)

// PayPal 配置
const (
	ClientID     = "AdE-SBHrFUjjnmwHEUJ5qzOGbzXdfQKaO-YAGfOW0R2XfzJAhWL4-iWaReNAkcqgtGb7OFnyIKklNNKA"
	ClientSecret = "EA_5m4fK_cizlnSrigdiQPTaWD1kx8x7sudb6rbs5T1gsStAhmYsInTeB_eqTegSi5cqNjDtxFjOaGhd"
	PayPalAPI    = "https://api-m.sandbox.paypal.com" // 沙箱环境
)

type OrderResponse struct {
	ID string `json:"id"`
}

// 获取 PayPal 访问令牌
func getAccessToken() (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", PayPalAPI+"/v1/oauth2/token", bytes.NewBufferString("grant_type=client_credentials"))
	if err != nil {
		return "", err
	}
	req.SetBasicAuth(ClientID, ClientSecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	return result["access_token"].(string), nil
}

// 创建订单
func CreatePayOrder(c echo.Context) error {
	accessToken, err := getAccessToken()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to get access token"})
	}

	order := map[string]interface{}{
		"intent": "CAPTURE",
		"purchase_units": []map[string]interface{}{
			{
				"amount": map[string]interface{}{
					"currency_code": "USD",
					"value":         "1.00",
				},
			},
		},
	}

	orderBody, _ := json.Marshal(order)

	client := &http.Client{}
	req, err := http.NewRequest("POST", PayPalAPI+"/v2/checkout/orders", bytes.NewBuffer(orderBody))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to create order"})
	}
	defer resp.Body.Close()

	var orderResponse OrderResponse
	if err := json.NewDecoder(resp.Body).Decode(&orderResponse); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to parse order response"})
	}

	return c.JSON(http.StatusOK, orderResponse)
}

// 捕获订单
func CaptureOrder(c echo.Context) error {
	orderID := c.Param("id")

	accessToken, err := getAccessToken()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to get access token"})
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", PayPalAPI+"/v2/checkout/orders/"+orderID+"/capture", nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to capture order"})
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	return c.JSON(http.StatusOK, json.RawMessage(body))
}
