package services

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

// プレイヤー情報
type Player struct {
	ID   string
	HP   int
	Conn *websocket.Conn
}

// デバイス情報
type Device struct {
	ID   string
	Conn *websocket.Conn
}

// ゲーム状態
type GameState struct {
	Player1HP int `json:"player1Hp"`
	Player2HP int `json:"player2Hp"`
}

var (
	devices   = map[string]*Device{} // デバイス情報を管理
	players   = map[string]*Player{} // プレイヤー情報を管理
	mu        sync.Mutex             // 同時アクセスを制御
	gameState = GameState{
		Player1HP: 100,
		Player2HP: 100,
	}
)

// デバイスを登録
func RegisterDevice(id string, conn *websocket.Conn) {
	mu.Lock()
	defer mu.Unlock()
	devices[id] = &Device{ID: id, Conn: conn}
	log.Printf("Device %s connected", id)
}

// デバイスの登録解除
func UnregisterDevice(id string) {
	mu.Lock()
	defer mu.Unlock()
	if device, exists := devices[id]; exists {
		device.Conn.Close()
		delete(devices, id)
		log.Printf("Device %s disconnected", id)
	}
}

// プレイヤーを登録
func RegisterPlayer(id string, conn *websocket.Conn) {
	mu.Lock()
	defer mu.Unlock()
	players[id] = &Player{ID: id, HP: 100, Conn: conn}
	log.Printf("Player %s connected", id)
}

// デバイスからの入力を処理
func ProcessInputFromDevice(deviceID string, message []byte) {
	mu.Lock()
	defer mu.Unlock()

	// 入力データを解析
	var input struct {
		Action string `json:"action"`
		Target string `json:"target"`
		Damage int    `json:"damage"`
	}
	if err := json.Unmarshal(message, &input); err != nil {
		log.Printf("Invalid input from device %s: %v", deviceID, err)
		return
	}

	// ゲームロジックの例: 攻撃処理
	if input.Target == "player1" {
		gameState.Player1HP -= input.Damage
		if gameState.Player1HP < 0 {
			gameState.Player1HP = 0
		}
	} else if input.Target == "player2" {
		gameState.Player2HP -= input.Damage
		if gameState.Player2HP < 0 {
			gameState.Player2HP = 0
		}
	}

	// 状態をブロードキャスト
	broadcastGameState()
}

// ゲーム状態を全プレイヤーに送信
func broadcastGameState() {
	for _, player := range players {
		if player.Conn != nil {
			err := player.Conn.WriteJSON(gameState)
			if err != nil {
				log.Printf("Error sending game state to player %s: %v", player.ID, err)
			}
		}
	}
}
