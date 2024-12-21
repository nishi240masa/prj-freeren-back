package models


// ゲーム状態
type GameState struct {
	Player1HP     int    `json:"player1Hp"`
	Player1MP     int    `json:"player1Mp"`
	Player1DF     int    `json:"player1Df"`
	Player1Action string `json:"player1Action"`
	Player2HP     int    `json:"player2Hp"`
	Player2MP     int    `json:"player2Mp"`
	Player2DF     int    `json:"player2Df"`
	Player2Action string `json:"player2Action"`
}