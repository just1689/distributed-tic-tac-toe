package model

type Player struct {
	ID               string `json:"id"`
	WebsocketChannel string `json:"channel"`
	Name             string `json:"name"`
}
