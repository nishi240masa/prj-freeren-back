package router

import (
	"fmt"
	"md2s/controllers"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Init() {
	// 環境変数 PORT を取得
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // デフォルトポート
	}

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},                                       // ここで特定のオリジンを許可することもできます（例: []string{"http://localhost:3000"})
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, // 許可するHTTPメソッド
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"}, // 許可するヘッダー
		AllowCredentials: true,                                                // クレデンシャルを許可するかどうか
		MaxAge:           12 * time.Hour,                                      // キャッシュの最大時間
	}))

	r.GET("/", func(c *gin.Context) {

		c.JSON(200, gin.H{
			"message": "connection success",
		})

	})

	// websocketのエンドポイント
	r.GET("/ws", controllers.HandleWebSocket)


	// デバイスの入力を処理するエンドポイント
	r.POST("/device/input", controllers.ProcessDeviceInputHandler)

	// 現在のゲーム状態を取得するエンドポイント
	r.GET("/game/state", controllers.GetGameStateHandler)

	// プレイヤーのWebSocket接続を処理するエンドポイント
	r.GET("/player/ws", controllers.HandlePlayerWebSocket)

	// 指定されたポートでサーバーを開始
	if err := r.Run(fmt.Sprintf(":%s", port)); err != nil {
		fmt.Printf("Failed to start server: %s\n", err)
	}
}
