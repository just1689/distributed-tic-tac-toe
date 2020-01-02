package server

import (
	"github.com/just1689/swoq/swoq"
)

var IncomingQueueName = "incoming"

func StartWS(listAddr string) {
	swoq.StartWebServer(listAddr, "/ws", IncomingQueueName)

}
