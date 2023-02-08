package scenarios_test

import (
	"sync"
	"testing"
	mockInteractor "zarg/lib/controller/mock"
	"zarg/lib/model/session"
)

func TestSingleRoom(t *testing.T) {
	interactor := mockInteractor.New(t)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		interactor.ExpectPrint() // greatings

		interactor.SendMsg("кот", "я")
		interactor.ExpectPrint() // кот в списке

		interactor.SendMsg("пес", "я")
		interactor.ExpectPrint() // пес в списке

		interactor.ExpectPrint() // 10 seconds until...
		interactor.ExpectPrint() // отправляемся в поход!

		interactor.ExpectPrint() // выбор оружия

		interactor.SendMsg("кот", "1")
		interactor.ExpectPrint() // кот выбрал оружие

		interactor.SendMsg("пес", "2")
		interactor.ExpectPrint() // пес выбрал оружие

		interactor.ExpectPrint() // выдвигаемся!

		interactor.ExpectPrint() // определите очерёдность

		interactor.SendMsg("кот", "я")
		interactor.ExpectPrint() // кот -> ...

		interactor.SendMsg("пес", "потом я")
		interactor.ExpectPrint() // кот -> пес

		interactor.EndInteraction()
		wg.Done()
	}()

	s := session.NewSession(interactor, func() {})
	wg.Wait()
	s.Stop()
}
