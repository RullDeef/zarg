package session

import (
	"context"
	"zarg/lib/model/enemy/squad"
	"zarg/lib/model/floormaze"
	I "zarg/lib/model/interfaces"
)

func (s *Session) exploreBossRoom(ctx context.Context, room *floormaze.BossRoom) {
	s.interactor.Printf("Вы находите просторную комнату, заваленную драгоценностями. Внезапно пол содрагается...")
	s.interactor.Printf("Перед вами появляется %s!", room.Boss.Name())

	es := squad.New(1, func() I.Enemy { return room.Boss })
	s.PerformBattle(ctx, es)
}
