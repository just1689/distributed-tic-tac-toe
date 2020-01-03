package model

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/just1689/swoq/queue"
	"github.com/sirupsen/logrus"
	"time"
)

var pingTime = time.Second * 2
var expiryTime = time.Second * 5

func NewServer(everyInstance, onlyOnce string) (s *Server) {
	s = &Server{
		ID:                    uuid.New().String(),
		IncomingEveryInstance: everyInstance,
		IncomingOnlyOnce:      onlyOnce,
		Players:               []*Player{},
		OtherInstances:        map[string]time.Time{},
		Games:                 []*Game{},
		chanAddPlayer:         make(chan *Player),
		chanRemovePlayer:      make(chan *Player),
		chanAddGame:           make(chan *Game),
		chanPublishAudit:      make(chan string),
		chanAddInstances:      make(chan string),
	}
	go s.handleChanges()
	go s.pingNetwork()
	return
}

func (s *Server) pingNetwork() {
	go func() {
		item := Message{
			Title: MessageIsInstanceID,
			Msg:   s.ID,
		}
		b, err := json.Marshal(item)
		if err != nil {
			logrus.Errorln("could not marshal WrappedMessage (pingNetwork)")
			logrus.Panicln(err)
		}
		publisher := queue.GetPublisher(s.IncomingEveryInstance)
		for {
			publisher(b)
			time.Sleep(pingTime)
		}
	}()
}

type Server struct {
	ID                    string               `json:"id"`
	IncomingEveryInstance string               `json:"-"`
	IncomingOnlyOnce      string               `json:"-"`
	Players               []*Player            `json:"players"`
	OtherInstances        map[string]time.Time `json:"otherInstances"`
	Games                 []*Game              `json:"games"`
	//Changes
	chanAddPlayer    chan *Player
	chanRemovePlayer chan *Player
	chanAddGame      chan *Game
	chanPublishAudit chan string
	chanAddInstances chan string
}

func (s *Server) HasPlayer(playerID string) bool {
	for _, p := range s.Players {
		if p.ID == playerID {
			return true
		}
	}
	return false
}
func (s *Server) GetPlayerByID(playerID string) (found bool, p *Player) {
	for _, p = range s.Players {
		if p.ID == playerID {
			found = true
			return
		}
	}
	return
}

func (s *Server) AddPlayer(p *Player) {
	s.chanAddPlayer <- p
}
func (s *Server) RemovePlayer(p *Player) {
	s.chanRemovePlayer <- p
}
func (s *Server) AddGame(g *Game) {
	s.chanAddGame <- g
}
func (s *Server) PublishAudit(channelID string) {
	s.chanPublishAudit <- channelID
}
func (s *Server) AddInstances(id string) {
	s.chanAddInstances <- id
}

func (s *Server) handleChanges() {
	for {
		select {
		case p := <-s.chanAddPlayer:
			logrus.Infoln("> Added player:", p.ID)
			s.Players = append(s.Players, p)
		case p := <-s.chanRemovePlayer:
			logrus.Infoln("> Removed player:", p.ID)
			s.Players = removePlayer(s.Players, p)
		case g := <-s.chanAddGame:
			logrus.Infoln("> Added game:", g.ID)
			s.Games = append(s.Games, g)
		case channelID := <-s.chanPublishAudit:
			logrus.Infoln("> publish audit")
			item := *s
			b, err := json.Marshal(item)
			if err != nil {
				logrus.Errorln(err)
				continue
			}
			m := Message{
				Title: MessageIsAuditResult,
				Msg:   "",
				Body:  b,
			}
			b, err = json.Marshal(m)
			if err != nil {
				logrus.Errorln(err)
				continue
			}
			queue.GetPublisher(channelID)(b)
		case instanceID := <-s.chanAddInstances:
			s.OtherInstances[instanceID] = time.Now()

			deadline := time.Now().Add(-1 * expiryTime)

			for key, lastSeen := range s.OtherInstances {
				if lastSeen.Before(deadline) {
					delete(s.OtherInstances, key)
					logrus.Info(key, " has not been seen for a while - removing")
				}
			}
		}
	}
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
