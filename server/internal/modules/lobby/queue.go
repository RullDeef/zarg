package lobby

import (
	"container/list"
	"errors"
	"server/domain"
	"sync"
	"time"
)

var ErrProfileAlreadyInQueue = errors.New("profile already in queue") // игрок уже находится в очереди

// requestQueue - очередь заявок на участие в походе одного типа
type requestQueue struct {
	queue        *list.List       // список структур partisipationRequest
	partyBuilder partyBuilderFunc // функция формирования команд
	mutex        sync.Mutex
	closed       bool
}

// participationRequest - запрос на участие в походе
type participationRequest struct {
	profile  *domain.Profile
	queuedAt time.Time
}

// partyBuilder - функция формирования команд.
// Возвращает true, если была успешно сформирована команда.
// Запросы участников команды при этом должны быть удалены из списка
type partyBuilderFunc func(*list.List) error

// newQueue - создает новую очередь заявок. функция partyBuilder
// начинает вызываться с периодом delay до момента закрытия очереди
func newQueue(partyBuilder partyBuilderFunc, delay time.Duration) *requestQueue {
	rq := requestQueue{
		queue:        list.New(),
		partyBuilder: partyBuilder,
		closed:       false,
	}

	go func() {
		for !rq.closed {
			<-time.After(delay)
			rq.checkParties()
		}
	}()

	return &rq
}

func (rq *requestQueue) checkParties() {
	rq.mutex.Lock()
	defer rq.mutex.Unlock()

	if rq.queue.Len() > 0 {
		if err := rq.partyBuilder(rq.queue); err != nil {
			panic(err)
		}
	}
}

// Close - завершает процесс формирования команд
func (rq *requestQueue) Close() error {
	rq.mutex.Lock()
	defer rq.mutex.Unlock()

	rq.closed = true
	return nil
}

// AddRequest - добавляет запрос на участие в походе в очередь
func (rq *requestQueue) AddRequest(profile *domain.Profile) error {
	rq.mutex.Lock()
	defer rq.mutex.Unlock()

	if rq.hasProfile(profile) {
		return ErrProfileAlreadyInQueue
	}

	rq.queue.PushBack(participationRequest{
		profile:  profile,
		queuedAt: time.Now(),
	})
	return nil
}

func (rq *requestQueue) hasProfile(profile *domain.Profile) bool {
	for node := rq.queue.Front(); node != nil; node = node.Next() {
		if req := node.Value.(participationRequest); req.profile == profile {
			return true
		}
	}
	return false
}

// CancelRequest - удаляет запрос из очереди
func (rq *requestQueue) CancelRequest(profile *domain.Profile) {
	rq.mutex.Lock()
	defer rq.mutex.Unlock()

	for node := rq.queue.Front(); node != nil; node = node.Next() {
		if req := node.Value.(participationRequest); req.profile == profile {
			rq.queue.Remove(node)
			return
		}
	}
}
