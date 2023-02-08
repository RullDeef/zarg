package scenarios_test

import (
	"sync"
	"testing"
	mockInteractor "zarg/lib/controller/mock"
	"zarg/lib/model/session"
)

func TestSinglePlayer(t *testing.T) {
	interactor := mockInteractor.New(t)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		interactor.ExpectPrint() // greatings

		interactor.SendMsg("кот", "я")

		interactor.ExpectPrint() // кот в списке
		interactor.ExpectPrint() // 10 seconds until...
		interactor.ExpectPrint() // одного недостаточно

		interactor.EndInteraction()
		wg.Done()
	}()

	s := session.NewSession(interactor, func() {})
	wg.Wait()
	s.Stop()
}
