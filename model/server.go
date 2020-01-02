package model

type Servers struct {
	list []Server
}

func (s *Servers) Add(item Server) {
	s.list = append(s.list, item)
}

type Server struct {
	players []Player
}
