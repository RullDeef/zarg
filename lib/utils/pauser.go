package utils

import (
	"context"
	"sync"
)

type Pauser struct {
	pausePub  *Publisher
	resumePub *Publisher
	isPaused  bool
	lock      sync.RWMutex
}

func NewPauser() *Pauser {
	return &Pauser{
		pausePub:  NewPublisher(),
		resumePub: NewPublisher(),
		isPaused:  false,
		lock:      sync.RWMutex{},
	}
}

func (p *Pauser) IsPaused() bool {
	p.lock.RLock()
	defer p.lock.RUnlock()
	return p.isPaused
}

func (p *Pauser) WaitPause(ctx context.Context) chan struct{} {
	p.lock.RLock()
	defer p.lock.RUnlock()

	c := make(chan struct{}, 1)

	if p.isPaused {
		c <- struct{}{}
	} else {
		go func() {
			s := NewSubscriber(p.pausePub)
			defer s.Unsubscribe()

			for {
				select {
				case <-s.Receive():
					c <- struct{}{}
					return
				case <-ctx.Done():
					return
				}
			}
		}()
	}

	return c
}

func (p *Pauser) WaitResume(ctx context.Context) chan struct{} {
	p.lock.RLock()
	defer p.lock.RUnlock()

	c := make(chan struct{}, 1)

	if !p.isPaused {
		c <- struct{}{}
	} else {
		go func() {
			s := NewSubscriber(p.resumePub)
			defer s.Unsubscribe()

			for {
				select {
				case <-s.Receive():
					c <- struct{}{}
					return
				case <-ctx.Done():
					return
				}
			}
		}()
	}

	return c
}

func (p *Pauser) Pause() {
	p.lock.Lock()
	defer p.lock.Unlock()

	if !p.isPaused {
		p.isPaused = true
		p.pausePub.Publish(struct{}{})
	}
}

func (p *Pauser) Resume() {
	p.lock.Lock()
	defer p.lock.Unlock()

	if p.isPaused {
		p.isPaused = false
		p.resumePub.Publish(struct{}{})
	}
}
