package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/just1689/distributed-tic-tac-toe/model"
	"github.com/just1689/distributed-tic-tac-toe/server"
	"github.com/just1689/swoq/queue"
	"github.com/just1689/swoq/swoq"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"

	"time"
)

var workers = flag.Int("workers", 12, "workers is the number of go routines for handling incoming requests")
var messageHandlers map[string]func(item *model.Message)

func main() {
	logrus.Println("Starting...")
	setupMsgHandlers()
	flag.Parse()
	if *workers <= 0 {
		logrus.Fatalln("Expected workers to be greater than 0, not ", *workers)
	}

	incomingWork := make(chan []byte, 1024)
	for i := 0; i < *workers; i++ {
		go startWorker(incomingWork)
		logrus.Println(" ...started worker ", i)
	}

	swoq.StartQueueClient()
	queueHandler := buildNATSHandler(incomingWork)

	logrus.Println("Starting subscriptions...")

	logrus.Println(" ...subscribing to", server.IncomingEveryInstance)
	queue.Subscribe(server.IncomingEveryInstance, queueHandler)

	logrus.Println(" ...subscribing to", server.IncomingOnlyOnce)
	if _, err := queue.DefaultConn.QueueSubscribe(server.IncomingOnlyOnce, "queue", queueHandler); err != nil {
		logrus.Fatalln(err)
	}

	go func() {
		time.Sleep(8 * time.Second)
		server.Instance.PublishAudit(server.IncomingOnlyOnce)
	}()

	logrus.Println("Backend has started!")
	select {}

}

func buildNATSHandler(incomingWork chan []byte) func(m *nats.Msg) {
	return func(m *nats.Msg) {
		incomingWork <- m.Data
	}
}

func startWorker(in chan []byte) {
	for b := range in {
		item := &model.Message{}
		err := json.Unmarshal(b, item)
		if err != nil {
			logrus.Errorln(err)
			continue
		}
		handleMessage(item)
	}
}

func setupMsgHandlers() {
	messageHandlers = make(map[string]func(item *model.Message))

	messageHandlers[model.MessageIsInstanceID] = func(item *model.Message) {
		server.Instance.AddInstances(item.Msg)
	}
	messageHandlers[model.MessageIsAuditResult] = func(item *model.Message) {
		fmt.Println(string(item.Body))
	}
	messageHandlers[model.MessageIsNewPlayer] = func(item *model.Message) {
		server.HandleNewPlayer(item)
	}
	messageHandlers[model.MessageIsNewGame] = func(item *model.Message) {
		server.NewGame(item)
	}
	messageHandlers[model.MessageIsGetPlayer] = func(item *model.Message) {
		server.HandleGetPlayerRemotely(item)
	}
	messageHandlers[model.MessageIsSetPlayer] = func(item *model.Message) {
		server.HandleSetPlayer(item)
	}
}

func handleMessage(item *model.Message) {
	f, found := messageHandlers[item.Title]
	if !found {
		fmt.Println("not sure how to handle", item.Title, item.Msg, string(item.Body))
		return
	}
	f(item)

}
