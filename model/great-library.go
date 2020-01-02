package model

type GreatLibrary struct {
	games []Game
}

func NewGreatLibrary() *GreatLibrary {
	return &GreatLibrary{
		games: make([]Game, 0),
	}
}
