package main

import (
	"flag"
	"github.com/just1689/distributed-tic-tac-toe/config"
	"github.com/just1689/distributed-tic-tac-toe/server"
	"github.com/just1689/swoq/queue"
	"github.com/sirupsen/logrus"
)

var listen = flag.String("listen", ":8080", "listen address")
var n = flag.String("nats", "nats://127.0.0.1:4222", "The NATS url for a NATS server instance.")

func main() {
	logrus.Println("Starting...")
	flag.Parse()
	queue.BuildDefaultConnFromUrl(*n)
	server.StartWS(config.GetVar("listen", *listen), server.IncomingOnlyOnce)

}
