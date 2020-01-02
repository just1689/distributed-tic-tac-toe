package server

import (
	"encoding/json"
	"github.com/just1689/swoq/queue"
	"github.com/just1689/swoq/ws"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
)

type Message struct {
	Name string `json:"name"`
}

func CreateQueueSubscriber() {
	queue.Subscribe(IncomingQueueName, func(m *nats.Msg) {
		w := &ws.WrappedMessage{}
		err := json.Unmarshal(m.Data, w)
		if err != nil {
			logrus.Errorln(err)
			return
		}
		msg := &Message{}
		err = json.Unmarshal(w.Body, msg)
		if err != nil {
			logrus.Errorln(err)
			return
		}
		reply := "hello " + msg.Name + "!"
		queue.Publish("client."+w.ClientID, reply)
	})

}
