package server

import (
	"encoding/json"
	"github.com/just1689/distributed-tic-tac-toe/model"
	"github.com/just1689/swoq/queue"
	"github.com/sirupsen/logrus"
	"time"
)

const IncomingEveryInstance = "global"
const IncomingOnlyOnce = "balanced"
const pingTime = time.Second * 2

var Instance *model.Server

func InitInstance() {
	Instance = model.NewServer(IncomingEveryInstance, IncomingOnlyOnce, []func(*model.Server){pingNetwork})
}

func pingNetwork(s *model.Server) {
	item := model.Message{
		Title: model.MessageIsInstanceID,
		Msg:   s.ID,
	}
	b, err := json.Marshal(item)
	if err != nil {
		logrus.Errorln("could not marshal WrappedMessage (pingNetwork)")
		logrus.Panicln(err)
	}
	publisher := queue.GetPublisher(s.IncomingEveryInstance)
	for {
		publisher(b)
		time.Sleep(pingTime)
	}
}
