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

var t = flag.Int("t", 0, "which test to run ")

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

	//Test A
	go func() {
		if *t != 1 {
			return
		}
		time.Sleep(3 * time.Second)
		p := &model.Player{
			ID:   "1000",
			Name: "Justin",
		}
		server.Instance.AddPlayer(p)
	}()
	//Test B
	go func() {
		if *t != 2 {
			return
		}
		time.Sleep(3 * time.Second)
		server.NewGame(&model.Message{
			Title: model.MessageIsNewGame,
			Msg:   "1000",
			Body:  nil,
		})
		server.Instance.PublishAudit(server.IncomingEveryInstance)
	}()

	go func() {
		time.Sleep(8 * time.Second)
		server.Instance.PublishAudit(server.IncomingEveryInstance)
	}()

	logrus.Println("Backend instance started", server.Instance.ID)
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

	messageHandlers[model.MessageIsInstanceID] = func(item *model.Message) { server.Instance.AddInstances(item.Msg) }
	messageHandlers[model.MessageIsAuditResult] = func(item *model.Message) { fmt.Println(string(item.Body)) }
	messageHandlers[model.MessageIsNewPlayer] = func(item *model.Message) { server.HandleNewPlayer(item) }
	messageHandlers[model.MessageIsNewGame] = func(item *model.Message) { server.NewGame(item) }
	messageHandlers[model.MessageIsGetPlayer] = func(item *model.Message) { server.HandleGetPlayerRemotely(item) }
	messageHandlers[model.MessageIsSetPlayer] = func(item *model.Message) { server.HandleSetPlayer(item) }
}

func handleMessage(item *model.Message) {
	f, found := messageHandlers[item.Title]
	if !found {
		fmt.Println("not sure how to handle", item.Title, item.Msg, string(item.Body))
		return
	}
	f(item)

}
