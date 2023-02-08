package mockInteractor

import (
	"context"
	"fmt"
	"testing"
	"zarg/lib/model"
	I "zarg/lib/model/interfaces"
	"zarg/lib/model/user"
)

type MockInteractor struct {
	receiver chan model.UserMessage
	printer  chan string
	users    map[string]*user.User
	t        *testing.T
}

func New(t *testing.T) *MockInteractor {
	return &MockInteractor{
		receiver: make(chan model.UserMessage, 1),
		printer:  make(chan string, 1),
		users:    make(map[string]*user.User),
		t:        t,
	}
}

// Interactor interface realisation
func (i *MockInteractor) Printf(format string, args ...any) {
	i.printer <- fmt.Sprintf(format, args...)
}

// Interactor interface realisation
func (i *MockInteractor) Receive(ctx context.Context, f func(I.UserMessage)) error {
	for {
		select {
		case umsg, ok := <-i.receiver:
			if ok {
				i.t.Logf("receiving [%s]: %s", umsg.User().FullName(), umsg.Message())
				f(umsg)
			} else {
				i.t.Error("tried to read from closed receiver chan!")
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (i *MockInteractor) SendMsg(userName string, msg string) {
	u := i.users[userName]
	if u == nil {
		i.users[userName] = user.New(len(i.users), userName, "")
		u = i.users[userName]
	}

	i.t.Logf("sending [%s]: %s", userName, msg)
	i.receiver <- model.NewUserMessage(u, msg)
}

func (i *MockInteractor) ExpectPrint() {
	msg := <-i.printer
	i.t.Logf("expected print: %s", msg)
}

func (i *MockInteractor) EndInteraction() {
	i.t.Log("ending interaction...")
	close(i.receiver)
}
