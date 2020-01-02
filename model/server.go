package model

func NewServer() *Server {
	return &Server{Players: []Player{}}
}

type Server struct {
	ID             string   `json:"id"`
	Players        []Player `json:"players"`
	OtherInstances []string `json:"otherInstances"`
}

func (s *Server) AddPlayer(p Player) {
	s.Players = append(s.Players, p)
}
