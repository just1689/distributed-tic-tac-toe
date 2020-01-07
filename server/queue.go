package server

import (
	"github.com/just1689/distributed-tic-tac-toe/model"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
)

func UnSubscribeFromQueues(hub *model.QueueHub) {
	hub.ForEach(func(name string, subscription *nats.Subscription) {
		logrus.Infoln("UnSubscribing from NATS queue:", name)
		err := subscription.Unsubscribe()
		if err != nil {
			logrus.Infoln("> FAIL")
			logrus.Errorln(err)
			return
		}
		logrus.Infoln("> OK ")

	})
}
