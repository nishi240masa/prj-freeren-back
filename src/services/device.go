package services

import (
	"encoding/json"
	"errors"
	"log"
	"md2s/models"
	"time"
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
func HttpProcessInputFromDevice(deviceID, action, state string) error {
	mu.Lock()
	defer mu.Unlock()

	// デバイスIDに基づいてプレイヤーを判定
	attacker, target := getPlayersByDevice(deviceID)
	if attacker == nil || target == nil {
		return errors.New("invalid device ID")
	}

	// stateを更新
	attacker.State = state

	if attacker.State == "noReady" || attacker.State == "" {
		log.Printf("Player %s is not ready", attacker.ID)
		return errors.New("player not ready")
	}

	// 片方が準備中の場合、actionを無視する
	if attacker.State == "ready" && target.State == "noReady" || target.State == "" {
		log.Printf("Player %s is ready, but Player %s is not ready", attacker.ID, target.ID)
		updateGameState()
		return errors.New("opponent not ready")
		
	}


	startCountdown := func() {

		if attacker.Time != 0 || target.Time != 0 {

		// カウントダウンを開始
			for i := 3; i > 0; i-- {
				attacker.Time = i
				target.Time = i
				updateGameState()
				time.Sleep(time.Second)
			}
			attacker.Time = 0
			target.Time = 0
			updateGameState()
		}
	}



	// 準備が完了した場合かつ相手も準備が完了している場合、カウントダウンを開始
	if attacker.State == "ready" && target.State == "ready" {
		startCountdown()

		return errors.New("change fighting")

	}

	// ゲームオーバーの場合、ゲームを終了
	if target.HP == 0 {
		log.Printf("Player %s wins!", attacker.ID)
		target.State = "death"
		attacker.State = "win"


		// 初期化
		attacker.HP = 100
		attacker.MP = 100
		attacker.DF = 100
		attacker.Action = "none"
		attacker.Time = 3

		target.HP = 100
		target.MP = 100
		target.DF = 100
		target.Action = "none"
		target.Time = 3



		updateGameState()



		return errors.New("game over")

	}



	if attacker.State == "fighting" && target.State == "fighting" {

	updateGameState()
		
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

	return errors.New("fighting")
}

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
		Player1State: players["player1"].State,
		Player2HP:     players["player2"].HP,
		Player2MP:     players["player2"].MP,
		Player2DF:     players["player2"].DF,
		Player2Action: players["player2"].Action,
		Player2State: players["player2"].State,
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
		damage := attacker.MP/5
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
		Player1State: players["player1"].State,
		Player2HP:     players["player2"].HP,
		Player2MP:     players["player2"].MP,
		Player2DF:     players["player2"].DF,
		Player2Action: players["player2"].Action,
		Player2State: players["player2"].State,
		Time:          players["player1"].Time,
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
