package session

import (
	"context"
	"fmt"
	"strings"
	"time"
	"zarg/lib/model/floormaze"
	I "zarg/lib/model/interfaces"
)

func (s *Session) exploreTrapRoom(ctx context.Context, room *floormaze.TrapRoom) {
	s.interactor.Printf("Что-то подсказывает вам, что тут не все так безобидно, как кажется на первый взгляд.")
	if s.makePauseFor(ctx, 4*time.Second) != nil {
		return
	}

	s.interactor.Printf(room.Trap.Name())
	if s.makePauseFor(ctx, 2*time.Second) != nil {
		return
	}

	healths := make(map[int]int)
	s.players.ForEachAlive(func(p I.Player) {
		healths[p.ID()] = p.Health()
	})
	damaged := room.Trap.Activate(s.players)
	info := ""

	var killedNames []string
	for _, p := range damaged {
		if !p.Alive() {
			killedNames = append(killedNames, p.FullName())
		} else {
			healthBefore := healths[p.ID()]
			info += fmt.Sprintf("%s (HP:%d->%d)\n", p.FullName(), healthBefore, p.Health())
		}
	}

	if killedNames != nil {
		if len(killedNames) == 1 {
			info += killedNames[0] + " убит!"
		} else {
			info += fmt.Sprintf("%s убиты!", strings.Join(killedNames, ", "))
		}
	}

	s.interactor.Printf(info)
}
