package model

import (
	"github.com/nats-io/nats.go"
)

func NewQueueHub() *QueueHub {
	return &QueueHub{
		list: make(map[string]*nats.Subscription),
	}
}

type QueueHub struct {
	list map[string]*nats.Subscription
}

func (q *QueueHub) Add(name string, sub *nats.Subscription) {
	q.list[name] = sub
}

func (q *QueueHub) Get(name string) (found bool, sub *nats.Subscription) {
	sub, found = q.list[name]
	return
}

func (q *QueueHub) ForEach(f func(name string, subscription *nats.Subscription)) {
	for name, subscription := range q.list {
		f(name, subscription)
	}
}
