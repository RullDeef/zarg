package session

import (
	"context"
	"strings"
	"time"
	"zarg/lib/model/floormaze"
	I "zarg/lib/model/interfaces"
)

func (s *Session) exploreRestRoom(ctx context.Context, room *floormaze.RestRoom) {
	info := "Вы находите комнату, в которой можно перевести дух и обговорить дальнейшие планы.\n"
	info += "Голосуйте \"в путь\" за то, чтобы продолжить поход, и \"строй\" за то, чтобы изменить очередность."
	s.Printf(info)

	if s.makePauseFor(ctx, 2*time.Second) != nil {
		return
	}

	s.players.ForEachAlive(func(p I.Entity) {
		p.Heal(50)
	})
	s.Printf("+50HP всем игрокам.")

	continueCounter := 0
	reorderCounter := 0
	reordering := false

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	s.receivePauseAware(ctx, func(umsg I.UserMessage) {
		p := s.players.GetByID(umsg.User().ID())
		if p == nil || reordering {
			return
		}

		msg := strings.ToLower(strings.TrimSpace(umsg.Message()))
		if msg == "в путь" {
			continueCounter += 1
			if 2*continueCounter > s.players.LenAlive() {
				cancel()
			}
		} else {
			continueCounter = 0
		}

		if msg == "строй" {
			reorderCounter += 1
			if 2*reorderCounter > s.players.LenAlive() {
				reorderCounter = 0
				reordering = true
				go func(ctx context.Context) {
					if !s.determinePlayersOrder(ctx) {
						s.Printf("Очередность не изменена.")
					}
					reordering = false
				}(ctx)
			}
		} else {
			reorderCounter = 0
		}
	})

	if ctx.Err() == nil {
		s.Printf("Поход продолжается!")
	}
}
