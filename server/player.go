package server

import (
	"github.com/google/uuid"
	"github.com/just1689/distributed-tic-tac-toe/model"
)

func HandleNewPlayer(message *model.Message) {
	p := &model.Player{
		ID:               uuid.New().String(),
		WebsocketChannel: message.Msg,
		Name:             "Unknown",
	}
	Instance.AddPlayer(p)
}
