package model

import (
	"github.com/google/uuid"
)

type Game struct {
	ID         string     `json:"id"`
	Players    []Player   `json:"players"`
	Board      [][]string `json:"board"`
	TurnType   TurnType   `json:"turnType"`
	PlayerTurn string     `json:"playerTurn"`
}

type TurnType string

const (
	PreGame  TurnType = "PreGame"
	InGame   TurnType = "InGame"
	PostGame TurnType = "PostGame"
)

func NewGame() *Game {
	result := &Game{
		ID:      uuid.New().String(),
		Players: make([]Player, 0),
		Board: [][]string{
			{"", "", ""},
			{"", "", ""},
			{"", "", ""},
		},
		TurnType:   PreGame,
		PlayerTurn: "",
	}
	return result
}

func (g *Game) HasPlayer(playerID string) bool {
	for _, p := range g.Players {
		if p.ID == playerID {
			return true
		}
	}
	return false
}

func (g *Game) GetPlayerByID(playerID string) (found bool, player Player) {
	for _, player = range g.Players {
		if player.ID == playerID {
			found = true
			return
		}
	}
	return
}

func (g *Game) GetPlayerChar(playerID string) string {
	if g.Players[0].ID == playerID {
		return "X"
	}
	return "0"
}
func (g *Game) GetCharPlayer(char string) string {
	if char == "X" {
		return g.Players[0].ID
	}
	return g.Players[1].ID
}

func (g *Game) PlayerMoves(playerID string, x, y int) (ok bool) {
	if !g.HasPlayer(playerID) {
		ok = false
		return
	}
	if x > 2 || x < 0 || y > 2 || y < 0 {
		ok = false
		return
	}
	if g.Board[x][y] != "" {
		ok = false
		return
	}

	g.Board[x][y] = g.GetPlayerChar(playerID)
	hasWinner := g.CheckForWinner()
	if hasWinner {
		g.TurnType = PostGame
	}
	ok = true
	return

}

func (g *Game) CheckForWinner() (hasWinner bool) {
	hasWinner = g.checkThreeForWinner(0, 0, 0, 1, 0, 2) ||
		g.checkThreeForWinner(1, 0, 1, 1, 1, 2) ||
		g.checkThreeForWinner(2, 0, 2, 1, 2, 2) ||
		g.checkThreeForWinner(0, 0, 1, 0, 2, 0) ||
		g.checkThreeForWinner(0, 1, 1, 1, 2, 1) ||
		g.checkThreeForWinner(0, 2, 1, 2, 2, 2) ||
		g.checkThreeForWinner(0, 0, 1, 1, 2, 2) ||
		g.checkThreeForWinner(0, 2, 1, 1, 2, 0)

	if hasWinner {
		g.TurnType = PostGame
	}
	return hasWinner

}

func (g *Game) checkThreeForWinner(x1, y1, x2, y2, x3, y3 int) bool {
	return g.Board[x1][y1] != "" && g.Board[x1][y1] == g.Board[x2][y2] && g.Board[x2][y2] == g.Board[x3][y3]

}

func (g *Game) StartGame() (ok bool) {
	if len(g.Players) != 2 {
		ok = false
		return
	}
	if g.TurnType != PreGame {
		ok = false
		return
	}
	g.TurnType = InGame
	ok = true
	return
}

func (g *Game) AddPlayer(p Player) (ok bool) {
	if len(g.Players) >= 2 {
		ok = false
		return
	}
	g.Players = append(g.Players, p)
	ok = true
	return
}
