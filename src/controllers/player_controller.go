package controllers

import (
	"log"
	"md2s/services"

	"github.com/gin-gonic/gin"
)

// HandlePlayerWebSocket プレイヤーのWebSocket接続を処理
func HandlePlayerWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}
	defer conn.Close()

	// クエリパラメータでプレイヤーIDを取得
	playerID := c.Query("player")
	if playerID == "" {
		log.Printf("No player ID provided")
		return
	}

	// プレイヤーを登録
	services.RegisterPlayer(playerID, conn)

	// メッセージの受信
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message from player %s: %v", playerID, err)
			break
		}

		// 入力を処理
		services.ProcessInputFromPlayer(playerID, message)
	}

}
