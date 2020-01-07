package model

import (
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
)

func NewQueueHub() *QueueHub {
	return &QueueHub{
		list: make(map[string]*nats.Subscription),
	}
}

type QueueHub struct {
	list map[string]*nats.Subscription
}

func (q *QueueHub) Add(name string, sub *nats.Subscription) {
	q.list[name] = sub
}

func (q *QueueHub) Get(name string) (found bool, sub *nats.Subscription) {
	sub, found = q.list[name]
	return
}

func (q *QueueHub) UnSubscribeAll() {
	for name, sub := range q.list {
		logrus.Infoln("Unsubscribing from NATS queue: ", name)
		err := sub.Unsubscribe()
		if err != nil {
			logrus.Infoln("> FAIL")
			logrus.Errorln(err)
			continue
		}
		logrus.Infoln("> OK ")
	}
}
