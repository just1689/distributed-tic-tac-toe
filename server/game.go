package server

import (
	"encoding/json"
	"github.com/just1689/distributed-tic-tac-toe/model"
	"github.com/sirupsen/logrus"
)

type newGameRequest struct {
	PlayerID string `json:"playerID"`
}

func NewGame(message *model.Message) {
	n := &newGameRequest{}
	err := json.Unmarshal(message.Body, n)
	if err != nil {
		logrus.Errorln("could not unmarshal newGameRequest")
		return
	}
	game := model.NewGame()

	if found, p := Instance.GetPlayer(n.PlayerID); found {
		game.AddPlayer(p)
	} else {
		game.AddPlayer(&model.Player{
			ID:   n.PlayerID,
			Name: "...",
		})
		game.FetchPlayerRemotely(n.PlayerID)
	}
	Instance.AddGame(game)
}
