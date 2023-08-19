package lobby

import (
	"container/list"
	"context"
	"errors"
	"server/domain"
	"sync"
	"time"
)

var (
	ErrProfileAlreadyInQueue = errors.New("profile already in queue") // игрок уже находится в очереди
	ErrRequestCancelled      = errors.New("request cancelled")        // запрос на участие в походе отменен
)

// requestQueue - очередь заявок на участие в походе одного типа
type requestQueue struct {
	queue        *list.List       // список структур *partisipationRequest
	partyBuilder partyBuilderFunc // функция формирования команд
	mutex        sync.Mutex
	closed       bool
}

// participationRequest - запрос на участие в походе
type participationRequest struct {
	profile    *domain.Profile
	queuedAt   time.Time
	joinSignal chan *party // специальный канал для определения, что игрок присоединился
}

// partyBuilder - функция формирования команд. Если команда была успешно сформирована,
// возвращается ее идентификатор, а также список участников, отобранных в команду.
// Передаваемый список не должен модифицироваться.
// Если команда пока еще не может быть сформирована, возвращается nil, nil
type partyBuilderFunc func(*list.List) (*party, error)

type party struct {
	compaignID domain.CompaignID
	requests   []*participationRequest
}

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
		if party, err := rq.partyBuilder(rq.queue); err != nil {
			panic(err)
		} else if party != nil {
			// check for join waiters here
			for _, waiter := range party.requests {
				waiter.joinSignal <- party
				close(waiter.joinSignal)
			}
			rq.removeRequests(party.requests)
		}
	}
}

// removeRequests - удаляет запросы из очереди
func (rq *requestQueue) removeRequests(requests []*participationRequest) {
	for _, req := range requests {
		for node := rq.queue.Front(); node != nil; node = node.Next() {
			if req1 := node.Value.(*participationRequest); req == req1 {
				rq.queue.Remove(node)
				break
			}
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

	rq.queue.PushBack(&participationRequest{
		profile:    profile,
		queuedAt:   time.Now(),
		joinSignal: make(chan *party, 1),
	})
	return nil
}

func (rq *requestQueue) hasProfile(profile *domain.Profile) bool {
	for node := rq.queue.Front(); node != nil; node = node.Next() {
		if req := node.Value.(*participationRequest); req.profile == profile {
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
		if req := node.Value.(*participationRequest); req.profile == profile {
			close(req.joinSignal)
			rq.queue.Remove(node)
			return
		}
	}
}

// CancelRequestByID - удаляет запрос из очереди по идентификатору профиля
func (rq *requestQueue) CancelRequestByID(profileID domain.ProfileID) {
	rq.mutex.Lock()
	defer rq.mutex.Unlock()

	for node := rq.queue.Front(); node != nil; node = node.Next() {
		if req := node.Value.(*participationRequest); req.profile.ID == profileID {
			close(req.joinSignal)
			rq.queue.Remove(node)
			return
		}
	}
}

// WaitJoin - ожидает формирования команды (блокирующий вызов).
// Если пользователь не находится ни в одной очереди - сразу возвращает ошибку.
func (rq *requestQueue) WaitJoin(ctx context.Context, profileID domain.ProfileID) (domain.CompaignID, error) {
	req := rq.getRequestByID(profileID)
	if req == nil {
		return "", ErrProfileNotInQueue
	}

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case party, ok := <-req.joinSignal:
		if !ok {
			return "", ErrRequestCancelled
		}
		return party.compaignID, nil
	}
}

func (rq *requestQueue) getRequestByID(id domain.ProfileID) *participationRequest {
	rq.mutex.Lock()
	defer rq.mutex.Unlock()

	for node := rq.queue.Front(); node != nil; node = node.Next() {
		if req := node.Value.(*participationRequest); req.profile.ID == id {
			return req
		}
	}
	return nil
}
