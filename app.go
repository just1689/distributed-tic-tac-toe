package main

import (
	"flag"
	"github.com/just1689/distributed-tic-tac-toe/server"
	"github.com/just1689/swoq/swoq"
	"github.com/sirupsen/logrus"
	"os"
)

var listen = flag.String("listen", "", "listen address")

func main() {
	flag.Parse()

	logrus.Println("Starting...")
	swoq.StartQueueClient()
	server.CreateQueueSubscriber()
	server.StartWS(getVar("listen"))

}

func getVar(key string) string {
	result := *listen
	if result == "" {
		result = os.Getenv(key)
	}
	return result
}
