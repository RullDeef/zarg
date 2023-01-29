package service

import (
	"log"
	"sync"
	"zarg/lib/model"
)

// service for managing interactors and sessions
type SessionBroker struct {
	interactors map[int]model.Interactor
	sessions    map[int]*model.Session
	lock        sync.RWMutex
}

func NewSessionBroker() *SessionBroker {
	return &SessionBroker{
		interactors: map[int]model.Interactor{},
		sessions:    map[int]*model.Session{},
		lock:        sync.RWMutex{},
	}
}

func (sb *SessionBroker) Interactor(chatID int, builder func() model.Interactor) model.Interactor {
	sb.lock.Lock()
	defer sb.lock.Unlock()

	i, ok := sb.interactors[chatID]
	if !ok {
		sb.interactors[chatID] = builder()
		i = sb.interactors[chatID]
	}
	return i
}

func (sb *SessionBroker) AddSession(chatID int, cleanup func()) bool {
	sb.lock.Lock()
	defer sb.lock.Unlock()

	if sb.sessions[chatID] != nil {
		return false
	}

	i := sb.interactors[chatID]
	if i == nil {
		log.Fatalf("failed to create session: interactor for chatID=%d is not set!", chatID)
	}

	sb.sessions[chatID] = model.NewSession(i, func() {
		defer cleanup()
		sb.lock.Lock()
		defer sb.lock.Unlock()

		sb.sessions[chatID] = nil
	})

	log.Printf("created session for chatID=%d\n", chatID)
	return true
}

func (sb *SessionBroker) Session(chatID int) *model.Session {
	sb.lock.Lock()
	defer sb.lock.Unlock()

	return sb.sessions[chatID]
}
