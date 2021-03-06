package server

import (
	"encoding/json"
	"github.com/just1689/distributed-tic-tac-toe/model"
	"github.com/just1689/swoq/queue"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
)

func NewGame(message model.Message) {
	game := model.NewGame()

	if found, p := Instance.GetPlayerByID(message.Msg); found {
		game.AddPlayer(p)
	} else {
		game.AddPlayer(&model.Player{
			ID:   message.Msg,
			Name: "...",
		})
		fetchPlayerRemotely(message.Msg, Instance.ID)
	}

	u := queue.Subscribe(game.GetChannel(), func(m *nats.Msg) {
		item := &model.Message{}
		if err := json.Unmarshal(m.Data, item); err != nil {
			logrus.Errorln(err)
			return
		}
		game.HandleIncomingMessage(item)
	})
	game.SubscriptionCloser = func() {
		u.Unsubscribe()
	}
	Instance.AddGame(game)
}

func fetchPlayerRemotely(playerID, backendID string) {
	m := model.Message{
		Title:  model.MessageIsGetPlayer,
		SrcKey: "backend",
		SrcID:  backendID,
		Msg:    playerID,
		Body:   nil,
	}
	b, err := json.Marshal(m)
	if err != nil {
		logrus.Errorln(err)
		return
	}
	queue.GetPublisher(Instance.IncomingEveryInstance)(b)
}
