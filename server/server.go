package server

import (
	"encoding/json"
	"github.com/just1689/distributed-tic-tac-toe/model"
	"github.com/just1689/swoq/queue"
	"github.com/sirupsen/logrus"
)

func HandleGetPlayerRemotely(item *model.Message) {
	found, p := Instance.GetPlayerByID(item.Msg)
	if found {
		item := *p
		b, err := json.Marshal(item)
		if err != nil {
			logrus.Errorln(err)
			return
		}
		m := model.Message{
			Title: model.MessageIsSetPlayer,
			Body:  b,
		}
		b, err = json.Marshal(m)
		if err != nil {
			logrus.Errorln(err)
			return
		}
		queue.GetPublisher(Instance.IncomingEveryInstance)(b)
	}
}

func HandleSetPlayer(item *model.Message) {
	p := &model.Player{}
	err := json.Unmarshal(item.Body, p)
	if err != nil {
		logrus.Errorln(err)
		return
	}
	Instance.AddPlayer(p)
}
