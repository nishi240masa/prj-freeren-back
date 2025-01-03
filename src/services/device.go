package services

import (
	"encoding/json"
	"errors"
	"log"
	"md2s/models"
)

// デバイス情報
type HttpDevice struct {
	ID string
}


// プレイヤーからの入力を処理
func ProcessInputFromPlayer(playerID string, message []byte) error {
	mu.Lock()
	defer mu.Unlock()

	player, exists := players[playerID]
	if !exists {
		return errors.New("player not found")
	}

	var input struct {
		Action string `json:"action"`
	}
	if err := json.Unmarshal(message, &input); err != nil {
		return err
	}

	// プレイヤーのアクションを更新

	player.Action = input.Action
	updateGameState()

	return nil
}


//  デバイスを登録
func HttpRegisterDevice(id string) error {
	mu.Lock()
	defer mu.Unlock()
	if _, exists := devices[id]; exists {
		return errors.New("device already registered")
	}
	devices[id] = &Device{ID: id}
	log.Printf("Device %s registered", id)
	return nil
}

//  デバイスの登録解除
func HttpUnregisterDevice(id string) {
	mu.Lock()
	defer mu.Unlock()
	delete(devices, id)
	log.Printf("Device %s unregistered", id)
}

// デバイスからの入力を処理
func HttpProcessInputFromDevice(deviceID, action string) error {
	mu.Lock()
	defer mu.Unlock()

	attacker, target := getPlayersByDevice(deviceID)
	if attacker == nil || target == nil {
		return errors.New("invalid device ID")
	}

	// プレイヤーの行動を登録
	attacker.Action = action

	// プレイヤーの行動を処理
	switch action {
	case "attack":
		processAttack(attacker, target)
	case "defend":
		processDefense(attacker)
	case "collection":
		processCollection(attacker)
	default:
		log.Printf("Unknown action: %s", action)
		return errors.New("unknown action")
	}

	// ゲーム状態を更新
	updateGameState()
	return nil
}

// GetGameState 現在のゲーム状態を取得
func GetGameState() models.GameState {
	mu.Lock()
	defer mu.Unlock()
	return models.GameState{
		Player1HP:     players["player1"].HP,
		Player1MP:     players["player1"].MP,
		Player1DF:     players["player1"].DF,
		Player1Action: players["player1"].Action,
		Player2HP:     players["player2"].HP,
		Player2MP:     players["player2"].MP,
		Player2DF:     players["player2"].DF,
		Player2Action: players["player2"].Action,
	}
}

// デバイスIDに対応するプレイヤーを取得
func getPlayersByDevice(deviceID string) (*Player, *Player) {
	if deviceID == "1" {
		return players["player1"], players["player2"]
	}
	if deviceID == "2" {
		return players["player2"], players["player1"]
	}
	return nil, nil
}

// 攻撃処理
func processAttack(attacker, target *Player) {
	if attacker.MP == 0 {
		log.Printf("Player %s has no MP", attacker.ID)
		return
	}
	if target.Action == "defend" {
		log.Printf("Player %s's attack was blocked by Player %s's defense!", attacker.ID, target.ID)
	} else {
		damage := attacker.MP / 10
		target.HP -= damage
		attacker.MP -= 20
		if attacker.MP < 0 {
			attacker.MP = 0
		}
		if target.HP < 0 {
			target.HP = 0
		}
		log.Printf("Player %s attacked Player %s for %d damage", attacker.ID, target.ID, damage)
	}
}

// 防御処理
func processDefense(player *Player) {

	if player.DF == 0 {
		log.Printf("Player %s has no DF", player.ID)
		player.Action = "none"
		return
	}

	player.DF -= 10
	if player.DF < 0 {
		player.DF = 0
	}
	log.Printf("Player %s is defending", player.ID)
}

// MP回復処理
func processCollection(player *Player) {
	player.MP += 10
	if player.MP > 100 {
		player.MP = 100
	}
	log.Printf("Player %s collected MP", player.ID)
}

// ゲーム状態を更新
func updateGameState() {
	gameState := models.GameState{
		Player1HP:     players["player1"].HP,
		Player1MP:     players["player1"].MP,
		Player1DF:     players["player1"].DF,
		Player1Action: players["player1"].Action,
		Player2HP:     players["player2"].HP,
		Player2MP:     players["player2"].MP,
		Player2DF:     players["player2"].DF,
		Player2Action: players["player2"].Action,
	}

	// プレイヤーにブロードキャスト
	for _, player := range players {
		if player.Conn != nil {
			err := player.Conn.WriteJSON(gameState)
			if err != nil {
				log.Printf("Error sending game state to player %s: %v", player.ID, err)
			}
		}
	}
}
