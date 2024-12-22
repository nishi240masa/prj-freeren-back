package controllers

import (
	"md2s/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

//  デバイスからの入力を処理
func ProcessDeviceInputHandler(c *gin.Context) {
	var input struct {
		DeviceID string `json:"deviceId"`
		Action   string `json:"action"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if err := services.HttpProcessInputFromDevice(input.DeviceID, input.Action); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Input processed successfully"})
}

//  現在のゲーム状態を取得
func GetGameStateHandler(c *gin.Context) {
	gameState := services.GetGameState()
	c.JSON(http.StatusOK, gameState)
}
