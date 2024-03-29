package session

import (
	"context"
	"fmt"
	"time"
	I "zarg/lib/model/interfaces"
	"zarg/lib/utils"
)

func (s *Session) Printf(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)

	s.logger.Printf("%s", msg)
	s.interactor.Printf("%s", msg)
}

// returns true if was canceled
func (s *Session) receiveWithAlert(ctx context.Context, d time.Duration, f func(umsg I.UserMessage, cancel func()), alertTime time.Duration, alertMsg string) bool {
	alarm := utils.AfterFunc(ctx, alertTime, s.pauser, func() {
		s.Printf("%s", alertMsg)
	})
	defer alarm.Stop()
	return s.receiveWithTimeout(ctx, d, f)
}

// returns true if was canceled
func (s *Session) receiveWithTimeout(ctx context.Context, d time.Duration, f func(umsg I.UserMessage, cancel func())) bool {
	ctx, cancel := s.timeoutFor(ctx, d)
	canceled := false
	s.receivePauseAware(ctx, func(umsg I.UserMessage) {
		f(umsg, func() {
			canceled = true
			cancel()
		})
	})
	return canceled
}

func (s *Session) receivePauseAware(ctx context.Context, f func(I.UserMessage)) error {
	return s.interactor.Receive(ctx, func(umsg I.UserMessage) {
		if !s.pauser.IsPaused() {
			s.logger.Printf("%s: %s", umsg.User().FullName(), umsg.Message())
			f(umsg)
		}
	})
}

func (s *Session) timeoutFor(parent context.Context, d time.Duration) (context.Context, func()) {
	ctx, cancel := context.WithCancel(parent)
	timer := utils.NewTimer(d, s.pauser)

	go func(parent context.Context) {
		defer cancel()
		<-timer.WaitCompleted(parent)
	}(parent)

	return ctx, cancel
}

func (s *Session) makePauseFor(ctx context.Context, d time.Duration) error {
	if ctx.Err() == nil {
		<-utils.NewTimer(d, s.pauser).WaitCompleted(ctx)
	}

	return ctx.Err()
}
