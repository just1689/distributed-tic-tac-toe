package main

import (
	"flag"
	"github.com/just1689/distributed-tic-tac-toe/config"
	"github.com/just1689/distributed-tic-tac-toe/model"
	"github.com/just1689/distributed-tic-tac-toe/server"
	"github.com/just1689/swoq/queue"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
	"time"
)

const NATsVar = "nats"

var workers = flag.Int("workers", 12, "workers is the number of go routines for handling incoming requests")
var natsURL = flag.String(NATsVar, "nats://127.0.0.1:4222", "The NATS url for a NATS server instance.")
var t = flag.Int("t", 0, "which test to run ")

func main() {
	logrus.Println("Starting...")
	server.SetupMsgHandlers()
	flag.Parse()
	if *workers <= 0 {
		logrus.Fatalln("Expected workers to be greater than 0, not ", *workers)
	}
	incomingWork := make(chan []byte, 1024)
	for i := 0; i < *workers; i++ {
		go server.StartWorker(incomingWork)
		logrus.Println(" ...started worker ", i)
	}
	natsConnURL := config.GetVar(NATsVar, *natsURL)
	logrus.Println("Connecting to NATs server @", natsConnURL)
	queue.BuildDefaultConnFromUrl(natsConnURL)

	logrus.Println("Starting subscriptions...")
	queueHandler := buildNATSHandler(incomingWork)
	logrus.Println(" ...subscribing to", server.IncomingEveryInstance)
	queue.Subscribe(server.IncomingEveryInstance, queueHandler)
	logrus.Println(" ...subscribing to", server.IncomingOnlyOnce)
	if _, err := queue.DefaultConn.QueueSubscribe(server.IncomingOnlyOnce, "queue", queueHandler); err != nil {
		logrus.Fatalln(err)
	}

	//FOR TEST
	if *t != 0 {
		setupTestEnv()
	}

	logrus.Println("Backend instance started", server.Instance.ID)
	select {}

}

func setupTestEnv() {
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
}

func buildNATSHandler(incomingWork chan []byte) func(m *nats.Msg) {
	return func(m *nats.Msg) {
		incomingWork <- m.Data
	}
}
