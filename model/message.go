package model

import "encoding/json"

var MessageIsInstanceID = "instance-id"
var MessageIsAuditResult = "audit-result"
var MessageIsNewPlayer = "new-player"
var MessageIsNewGame = "new-game"
var MessageIsGetPlayer = "get-player"
var MessageIsSetPlayer = "set-player"

type Message struct {
	Title  string          `json:"title"`
	SrcKey string          `json:"srcKey"`
	SrcID  string          `json:"srcID"`
	Msg    string          `json:"msg"`
	Body   json.RawMessage `json:"body"`
}

func (m *Message) GetReplyChannel() string {
	return m.SrcKey + "." + m.SrcID
}
