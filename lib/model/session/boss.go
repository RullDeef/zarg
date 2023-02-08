package session

import (
	"context"
	"zarg/lib/model"
	"zarg/lib/model/enemy/boss"
	"zarg/lib/model/enemy/squad"
	I "zarg/lib/model/interfaces"
	"zarg/lib/utils"
)

func (s *Session) exploreBossRoom(ctx context.Context, fm *model.FloorMaze) {
	s.interactor.Printf("Вы находите просторную комнату, заваленную драгоценностями. Внезапно пол содрагается...")

	bossesFns := make(map[string]func(context.Context, *model.FloorMaze))
	bossesFns["troll"] = s.troll

	boss := utils.NewPropMap()
	for k := range bossesFns {
		boss.Add(k, 1)
	}
	bossesFns[boss.Choose().(string)](ctx, fm)
}

func (s *Session) troll(ctx context.Context, fm *model.FloorMaze) {
	s.interactor.Printf("Перед вами появляется огромный тролль!")

	boss := boss.New(
		boss.NewPhase("Тролль", 300, func() I.DamageStats {
			return I.DamageStats{
				Base:       30,
				Crit:       60,
				CritChance: 0.25,
			}
		}, func(bp1, bp2 *boss.BossPhase) {
			s.interactor.Printf("Троль разгневался и стал сильнее!")
		}),
		boss.NewPhase("Разъяренный Тролль", 200, func() I.DamageStats {
			return I.DamageStats{
				Base:       50,
				Crit:       80,
				CritChance: 0.4,
			}
		}, nil),
	)

	es := squad.New(1, func() I.Enemy { return boss })
	s.PerformBattle(ctx, es)
}
