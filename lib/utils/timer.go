package utils

import (
	"context"
	"time"
)

type Timer struct {
	at       time.Time
	duration time.Duration
	pauser   *Pauser
}

type Timeouter struct {
	Timer
	canceled bool
}

func NewTimer(d time.Duration, pauser *Pauser) *Timer {
	return &Timer{
		at:       time.Now(),
		duration: d,
		pauser:   pauser,
	}
}

func AfterFunc(ctx context.Context, d time.Duration, pauser *Pauser, action func()) *Timeouter {
	t := &Timeouter{
		Timer: Timer{
			at:       time.Now(),
			duration: d,
			pauser:   pauser,
		},
		canceled: false,
	}

	go func(ctx context.Context) {
		<-t.WaitCompleted(ctx)
		if !t.canceled {
			action()
		}
	}(ctx)

	return t
}

func (t *Timer) WaitCompleted(ctx context.Context) chan struct{} {
	c := make(chan struct{}, 1)

	go func() {
	outer:
		for t.duration > 0 {
			select {
			case <-time.After(t.duration):
				break outer
			case <-t.pauser.WaitPause(ctx):
				t.duration -= time.Since(t.at)
				<-t.pauser.WaitResume(ctx)
				t.at = time.Now()
			}
		}
		c <- struct{}{}
	}()

	return c
}

func (t *Timeouter) Stop() {
	t.canceled = true
}
