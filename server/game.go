package server

import (
	"encoding/json"
	"github.com/just1689/distributed-tic-tac-toe/model"
	"github.com/just1689/swoq/queue"
	"github.com/sirupsen/logrus"
)

func NewGame(message *model.Message) {
	game := model.NewGame()

	if found, p := Instance.GetPlayerByID(message.Msg); found {
		game.AddPlayer(p)
	} else {
		game.AddPlayer(&model.Player{
			ID:   message.Msg,
			Name: "...",
		})
		fetchPlayerRemotely(message.Msg)
	}
	Instance.AddGame(game)
}

func fetchPlayerRemotely(id string) {
	m := model.Message{
		Title: model.MessageIsGetPlayer,
		Msg:   id,
	}
	b, err := json.Marshal(m)
	if err != nil {
		logrus.Errorln(err)
		return
	}
	queue.GetPublisher(Instance.IncomingEveryInstance)(b)
}
