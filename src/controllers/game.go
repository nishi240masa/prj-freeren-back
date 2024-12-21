package controllers

import (
	"log"
	"md2s/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// WebSocketのアップグレーダー
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// HandleWebSocket handles WebSocket connections
func HandleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}
	defer conn.Close()

	// クエリパラメータでデバイスIDを取得
	deviceID := c.Query("deviceId")
	playerId := c.Query("playerId")

	// もしどっちもない場合はエラー
	if deviceID == "" && playerId == "" {
		log.Printf("No device ID or player ID provided")
		return
	}

	// デバイスを登録
	services.RegisterDevice(deviceID, conn)

	// プレイヤーを登録
	services.RegisterPlayer(playerId, conn)

	// メッセージの受信
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message from device %s: %v", deviceID, err)
			break
		}

		// 入力を処理
		services.ProcessInputFromDevice(deviceID, message)
	}

	// デバイスの登録解除
	services.UnregisterDevice(deviceID)
}
