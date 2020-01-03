package server

import (
	"github.com/just1689/swoq/swoq"
)

func StartWS(listAddr, globalIncoming string) {
	swoq.StartWebServer(listAddr, "/ws", globalIncoming)

}
