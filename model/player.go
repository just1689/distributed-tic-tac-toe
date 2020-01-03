package model

type Player struct {
	ID               string `json:"id"`
	WebsocketChannel string `json:"channel"`
	Name             string `json:"name"`
}

func removePlayer(all []*Player, remove *Player) []*Player {
	var i int
	var p *Player
	for i, p = range all {
		if p.ID == remove.ID {
			break
		}
	}
	all[i] = all[len(all)-1]
	return all[:len(all)-1]
}
