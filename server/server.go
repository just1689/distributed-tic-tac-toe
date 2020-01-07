package server

import (
	"encoding/json"
	"fmt"
	"github.com/just1689/distributed-tic-tac-toe/model"
	"github.com/just1689/swoq/queue"
	"github.com/sirupsen/logrus"
)

var messageHandlers map[string]func(item *model.Message)

func SetupMsgHandlers() {
	messageHandlers = map[string]func(message *model.Message){
		model.MessageIsAuditResult: printBody,
		model.MessageIsInstanceID:  AddInstance,
		model.MessageIsNewPlayer:   HandleNewPlayer,
		model.MessageIsNewGame:     NewGame,
		model.MessageIsGetPlayer:   HandleGetPlayerRemotely,
		model.MessageIsSetPlayer:   HandleSetPlayer,
	}
}

func printBody(item *model.Message) {
	fmt.Println(string(item.Body))
}

func StartWorker(in chan []byte) {
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

func handleMessage(item *model.Message) {
	if f, found := messageHandlers[item.Title]; !found {
		fmt.Println("not sure how to handle", item.Title, item.Msg, string(item.Body))
		return
	} else {
		f(item)
	}
}

func AddInstance(item *model.Message) {
	Instance.AddInstances(item.Msg)
}

func HandleGetPlayerRemotely(item *model.Message) {
	found, p := Instance.GetPlayerByID(item.Msg)
	if found {
		b, err := json.Marshal(*p)
		if err != nil {
			logrus.Errorln(err)
			return
		}
		m := model.Message{
			Title: model.MessageIsSetPlayer,
			Body:  b,
		}
		b, err = json.Marshal(m)
		if err != nil {
			logrus.Errorln(err)
			return
		}
		queue.GetPublisher(item.GetReplyChannel())(b)
	}
}

func HandleSetPlayer(item *model.Message) {
	player := &model.Player{}
	err := json.Unmarshal(item.Body, player)
	if err != nil {
		logrus.Errorln(err)
		return
	}
	for _, g := range Instance.Games {
		for i, p := range g.Players {
			if p.ID == player.ID {
				g.Players[i] = player
				return
			}
		}
	}
}
