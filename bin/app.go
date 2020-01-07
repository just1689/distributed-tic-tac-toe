package main

import (
	"encoding/json"
	"flag"
	"github.com/just1689/distributed-tic-tac-toe/config"
	"github.com/just1689/distributed-tic-tac-toe/model"
	"github.com/just1689/distributed-tic-tac-toe/server"
	"github.com/just1689/swoq/queue"
	"github.com/just1689/swoq/ws"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

const NATsVar = "nats"

var app = flag.String("app", "", "Which app to run - backend or gateway")
var workers = flag.Int("workers", 12, "workers is the number of go routines for handling incoming requests")
var natsURL = flag.String(NATsVar, "nats://127.0.0.1:4222", "The NATS url for a NATS server instance.")
var t = flag.Int("t", 0, "which test to run ")
var listen = flag.String("listen", ":8080", "listen address")
var shutdownSeconds = flag.Int("shutdownSeconds", 4, "how many seconds the app has to shutdown")

func main() {
	flag.Parse()
	runSignalHandler()
	a := config.GetVar("app", *app)
	if a == "backend" {
		RunBackend()
	} else if a == "gateway" {
		RunGateway()
	}
	os.Exit(1)
	logrus.Errorln("could not start app", a)

}

func RunGateway() {
	logrus.Println("Starting...")
	flag.Parse()
	n := config.GetVar(NATsVar, *natsURL)
	queue.BuildDefaultConnFromUrl(n)
	server.StartWS(config.GetVar("listen", *listen), server.IncomingWebsocket)

}

func RunBackend() {
	logrus.Println("Starting...")
	server.InitInstance()
	server.SetupMsgHandlers()
	if *workers <= 0 {
		logrus.Fatalln("Expected workers to be greater than 0, not ", *workers)
	}
	incomingBin := make(chan []byte, 1024)
	incomingItem := make(chan model.Message, 1024)
	for i := 0; i < *workers; i++ {
		go server.StartConverter(incomingBin, incomingItem)
		go server.StartWorker(incomingItem)
		logrus.Println(" ...started worker ", i)
	}

	connectToNATS()
	subscribeToNATS(incomingBin, incomingItem)

	//FOR DISTRIBUTED TEST
	if *t != 0 {
		setupTestEnv()
	}

	logrus.Println("Backend instance started", server.Instance.ID)
	select {}
}

func subscribeToNATS(incomingBin chan []byte, incomingItem chan model.Message) {
	logrus.Println("Starting subscriptions...")
	queueHandler := buildNATSHandler(incomingBin)
	logrus.Println(" ...subscribing to", server.IncomingEveryInstance)
	unSubEveryInstance := queue.Subscribe(server.IncomingEveryInstance, queueHandler)
	server.Instance.QueueHub.Add(server.IncomingEveryInstance, unSubEveryInstance)
	logrus.Println(" ...subscribing to", server.IncomingOnlyOnce)
	unSubQueueWS, err := queue.DefaultConn.QueueSubscribe(server.IncomingWebsocket, "queue", buildWebsocketHandler(incomingItem))
	if err != nil {
		logrus.Fatalln(err)
	}
	server.Instance.QueueHub.Add(server.IncomingWebsocket, unSubQueueWS)
	logrus.Println(" ...subscribing to", server.Instance.GetQueueName())
	unSubQueueBackendQueue := queue.Subscribe(server.Instance.GetQueueName(), queueHandler)
	server.Instance.QueueHub.Add(server.Instance.GetQueueName(), unSubQueueBackendQueue)

}

func connectToNATS() {
	natsConnURL := config.GetVar(NATsVar, *natsURL)
	logrus.Println("Connecting to NATs server @", natsConnURL)
	queue.BuildDefaultConnFromUrl(natsConnURL)

}

func buildWebsocketHandler(incomingItem chan model.Message) func(m *nats.Msg) {
	return func(msg *nats.Msg) {
		item := &ws.WrappedMessage{}
		if err := json.Unmarshal(msg.Data, item); err != nil {
			logrus.Errorln("could not convert websocket message to WrappedMessage")
			logrus.Errorln(string(msg.Data))
			logrus.Errorln(err)
			return
		}
		result := model.Message{
			Title:  "incoming-ws",
			SrcKey: "client",
			SrcID:  item.ClientID,
			Msg:    "",
			Body:   item.Body,
		}
		incomingItem <- result
	}
}

func buildNATSHandler(incomingWork chan []byte) func(m *nats.Msg) {
	return func(m *nats.Msg) {
		incomingWork <- m.Data
	}
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
		server.NewGame(model.Message{
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

func runSignalHandler() {
	s := strconv.Itoa(*shutdownSeconds) + "s"
	d, err := time.ParseDuration(s)
	if err != nil {
		logrus.Errorln("could not understand shutdownSeconds: ", *shutdownSeconds, " as an int")
		logrus.Fatalln(err)
	}

	go func(d time.Duration) {
		// channel for incoming signals
		sigint := make(chan os.Signal, 1)
		// interrupt signal sent from terminal
		signal.Notify(sigint, os.Interrupt)
		// sigterm signal sent from kubernetes
		signal.Notify(sigint, syscall.SIGTERM)

		<-sigint
		logrus.Println("Terminate signal received, beginning shutdown")
		go server.UnSubscribeFromQueues(server.Instance.QueueHub)
		logrus.Println("Stopping in", d)
		time.Sleep(d)
		os.Exit(0)

	}(d)
}
