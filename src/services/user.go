package services

import (
	"encoding/json"
	"log"
	"md2s/models"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// プレイヤー情報
type Player struct {
	ID     string
	HP     int
	MP     int
	DF     int
	Action string  // 現在の行動 ("attack", "defend", etc.)
	// 準備中か戦闘中かなどの状態
	State string   // 現在の状態 ("noReady","ready", etc.)
	Time  int
	Conn   *websocket.Conn
}


// デバイス情報
type Device struct {
	ID   string
	Conn *websocket.Conn
}



var (
	devices   = map[string]*Device{} // デバイス情報を管理
	players   = map[string]*Player{} // プレイヤー情報を管理
	mu        sync.Mutex             // 同時アクセスを制御
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
	// playerの初期値を設定
    players[id] = &Player{ID: id, HP: 100, MP: 100, DF: 100, Action: "none",State: "noReady",Time: 3, Conn: conn}
	log.Printf("Player %s connected", id)

	// 状態をブロードキャスト
	broadcastGameState()
}

// デバイスからの入力を処理
func ProcessInputFromDevice(deviceID string, message []byte) {
	mu.Lock()
	defer mu.Unlock()

	// 入力データを解析
	var input struct {
		Action string `json:"action"`
		State string `json:"state"`
	}
	if err := json.Unmarshal(message, &input); err != nil {
		log.Printf("Invalid input from device %s: %v", deviceID, err)
		return
	}

	// デバイスIDに基づいてプレイヤーを判定
	var attackingPlayer *Player
	var targetPlayer *Player
	if deviceID == "1" {
		attackingPlayer = players["player1"]
		targetPlayer = players["player2"]
	} else if deviceID == "2" {
		attackingPlayer = players["player2"]
		targetPlayer = players["player1"]
	}

	if attackingPlayer == nil || targetPlayer == nil {
		log.Printf("Invalid device ID: %s", deviceID)
		return
	}

	// stateを更新
	attackingPlayer.State = input.State

	// 準備中の場合、actionを無視する

	if attackingPlayer.State == "noReady" {
		return
	}

	// 準備が完了した場合かつ相手も準備が完了している場合、戦闘を開始actionを処理
	if attackingPlayer.State == "ready" && targetPlayer.State == "ready" {


		// 3秒カウントダウン
		// 3秒後に攻撃処理を実行
		attackingPlayer.Time = 3
		targetPlayer.Time = 3

		// ゲーム状態をブロードキャスト
		broadcastGameState()

		// カウントダウン処理
		for attackingPlayer.Time > 0 {
			attackingPlayer.Time--
			targetPlayer.Time--
			broadcastGameState()
			time.Sleep(1 * time.Second)
		}


		// 3秒後にstateを戦闘中に変更
		attackingPlayer.State = "fighting"
		targetPlayer.State = "fighting"

		// ゲーム状態をブロードキャスト
		broadcastGameState()






	// プレイヤーの行動を更新
	attackingPlayer.Action = input.Action

	// ゲームロジック: 攻撃処理
	if input.Action == "attack" {
		// 相手が防御している場合、攻撃は無効
		if targetPlayer.Action == "defend" {
			log.Printf("Player %s's attack was blocked by Player %s's defense!", attackingPlayer.ID, targetPlayer.ID)
		} else {
			// 防御していない場合、ダメージを与える
			damage := attackingPlayer.MP / 10 // MPの10%をダメージとして与える
			targetPlayer.HP -= damage
			attackingPlayer.MP -= 20
			if targetPlayer.HP < 0 {
				targetPlayer.HP = 0
			}
		}
	}

    // ゲームロジック: 防御処理
    if input.Action == "defend" {
        attackingPlayer.DF -= 10
        if attackingPlayer.DF < 0 {
            attackingPlayer.DF = 0
        }
    }

    // ゲームロジック: MP回復
    if input.Action == "collection" {
        attackingPlayer.MP += 20
        if attackingPlayer.MP > 100 {
            attackingPlayer.MP = 100
        }
    }




	// ゲーム終了判定
	if targetPlayer.HP == 0 {
		log.Printf("Player %s wins!", attackingPlayer.ID)
				attackingPlayer.State = "noReady"
		targetPlayer.State = "noReady"

		// プレイヤーのHPとMP,DFをリセット
		attackingPlayer.HP = 100
		targetPlayer.HP = 100
		attackingPlayer.MP = 100
		targetPlayer.MP = 100
		attackingPlayer.DF = 100
		targetPlayer.DF = 100


	}



	// 状態をブロードキャスト
	broadcastGameState()
}

}


// ゲーム状態を全プレイヤーに送信
func broadcastGameState() {
	gameState := models.GameState{
		Player1HP:   players["player1"].HP,
		Player1MP:   players["player1"].MP,
		Player1DF:   players["player1"].DF,
		Player1Action: players["player1"].Action,
		Player1State: players["player1"].State,
		Player2HP:   players["player2"].HP,
		Player2MP:   players["player2"].MP,
		Player2DF:   players["player2"].DF,
		Player2Action: players["player2"].Action,
		Player2State: players["player2"].State,
	}

	for _, player := range players {
		if player.Conn != nil {
			err := player.Conn.WriteJSON(gameState)
			if err != nil {
				log.Printf("Error sending game state to player %s: %v", player.ID, err)
			}
		}
	}
}

