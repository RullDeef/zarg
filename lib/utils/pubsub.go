package utils

import (
	"log"
	"sync"
)

type Publisher struct {
	subs []*Subscriber
	lock sync.RWMutex
}

type Subscriber struct {
	pub  *Publisher
	id   int
	pipe chan any
}

func NewPublisher() *Publisher {
	return &Publisher{
		subs: nil,
	}
}

func NewSubscriber(p *Publisher) *Subscriber {
	p.lock.Lock()
	defer p.lock.Unlock()

	s := &Subscriber{
		pub:  p,
		id:   len(p.subs),
		pipe: make(chan any),
	}

	p.subs = append(p.subs, s)
	log.Printf("sub #%p subscribed", s)
	return s
}

func (p *Publisher) Publish(data any) {
	p.lock.RLock()
	defer p.lock.RUnlock()
	for _, s := range p.subs {
		s.pipe <- data
	}
}

func (p *Publisher) Close() {
	p.lock.Lock()
	defer p.lock.Unlock()

	for _, s := range p.subs {
		close(s.pipe)
	}
	p.subs = nil
}

func (s *Subscriber) Receive() chan any {
	return s.pipe
}

func (s *Subscriber) Unsubscribe() {
	s.pub.lock.Lock()
	defer s.pub.lock.Unlock()

	if s.id+1 == len(s.pub.subs) {
		s.pub.subs = s.pub.subs[:len(s.pub.subs)-1]
	} else {
		s.pub.subs[s.id] = s.pub.subs[len(s.pub.subs)-1]
		s.pub.subs = s.pub.subs[:len(s.pub.subs)-1]
		s.pub.subs[s.id].id = s.id
	}

	log.Printf("sub #%p unsubscribed", s)
}
