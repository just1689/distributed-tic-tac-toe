package server

import "github.com/just1689/distributed-tic-tac-toe/model"

var IncomingEveryInstance = "global"
var IncomingOnlyOnce = "balanced"

var Instance = model.NewServer(IncomingEveryInstance, IncomingOnlyOnce)
