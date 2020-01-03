package model

import "encoding/json"

var MessageIsInstanceID = "instance-id"
var MessageIsAuditResult = "audit-result"
var MessageIsNewPlayer = "new-player"

type Message struct {
	Title string          `json:"title"`
	Msg   string          `json:"msg"`
	Body  json.RawMessage `json:"body"`
}
