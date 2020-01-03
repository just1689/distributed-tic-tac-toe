package server

import "github.com/just1689/distributed-tic-tac-toe/model"

type newGameRequest struct {
	PlayerID string `json:"playerID"`
}

func NewGame(message model.Message) {
	game := model.NewGame()

	Instance.AddGame(game)
}
