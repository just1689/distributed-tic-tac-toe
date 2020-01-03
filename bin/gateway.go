package main

import (
	"flag"
	"github.com/just1689/distributed-tic-tac-toe/config"
	"github.com/just1689/distributed-tic-tac-toe/server"
	"github.com/just1689/swoq/swoq"
	"github.com/sirupsen/logrus"
)

var listen = flag.String("listen", "", "listen address")

func main() {
	logrus.Println("Starting...")
	swoq.StartQueueClient()
	server.StartWS(config.GetVar("listen", *listen), server.IncomingOnlyOnce)

}
