package session

import (
	"context"
	"fmt"
	"strconv"
	"time"
	enemySquad "zarg/lib/model/enemy/squad"
	"zarg/lib/model/floormaze"
	I "zarg/lib/model/interfaces"
	"zarg/lib/utils"
)

func (s *Session) exploreEnemiesRoom(ctx context.Context, room *floormaze.EnemyRoom) {
	s.Printf("Вы не одни... На вас напали!")
	s.PerformBattle(ctx, room.Enemies)
}

func (s *Session) PerformBattle(ctx context.Context, es *enemySquad.EnemySquad) {

	// show battle overall info
	battleInfo := fmt.Sprintf("Игроки:\n%sВраги:\n%s", s.players.CompactInfo(), es.CompactInfo())
	s.Printf(battleInfo)

	// TODO: make general turn referee
	turnsMadePlayers := 0
	turnsMadeEnemies := 0

	for es.LenAlive() > 0 && s.players.LenAlive() > 0 {
		if s.makePauseFor(ctx, time.Second) != nil {
			return
		}

		// make turn generator
		turnGen := utils.NewPropMap()
		turnGen.Add("players", utils.MaxInt(0, s.players.LenAlive()-turnsMadePlayers))
		turnGen.Add("enemies", utils.MaxInt(0, es.LenAlive()-turnsMadeEnemies))

		switch turnGen.Choose().(string) {
		case "players":
			// show battle overall info
			s.Printf("Игроки:\n%sВраги:\n%s", s.players.CompactInfo(), es.CompactInfo())
			p := s.players.ChooseNext()
			s.makePlayerAction(ctx, p, es)
			turnsMadePlayers += 1
		case "enemies":
			e := es.ChooseNext()
			s.makeEnemyAction(ctx, e, es)
			turnsMadeEnemies += 1
		default:
			s.logger.Panicf("must never happen!")
		}

		if turnsMadePlayers >= s.players.LenAlive() && turnsMadeEnemies >= es.LenAlive() {
			turnsMadePlayers = 0
			turnsMadeEnemies = 0
		}
	}

	if es.LenAlive() == 0 {
		s.Printf("Битва завершена. Все враги повержены!")
	} else {
		s.Printf("Битва завершена. Все игроки мертвы!")
	}
}

func (s *Session) makePlayerAction(ctx context.Context, p I.Player, es *enemySquad.EnemySquad) {
	s.Printf("Ходит %s.", p.FullName())

	if s.makePauseFor(ctx, 5*time.Second) != nil {
		return
	}

	optsInfo := "Варианты действия:\n"
	opts := map[int]func(){}
	i := 1
	es.ForEachAlive(func(e I.Enemy) {
		optsInfo += fmt.Sprintf("%d) Атаковать %s (HP: %d)\n", i, e.Name(), e.Health())
		opts[i] = func() {
			dmg := e.Damage(p.Attack())
			if e.Alive() {
				s.Printf("%s атакует %s и наносит %d урона.", p.FullName(), e.Name(), dmg)
			} else {
				s.Printf("%s убивает %s.", p.FullName(), e.Name())
			}
		}
		i++
	})
	// add block option
	optsInfo += fmt.Sprintf("%d) Поставить блок (x0.8DMG)\n", i)
	opts[i] = func() {
		p.BlockAttack()
		s.Printf("%s ставит блок!", p.FullName())
	}
	i++
	p.ForEachItem(func(item I.Pickable) {
		if usable, ok := item.(I.Usable); ok {
			if !usable.IsUsed() {
				opts[i] = func() {
					usable.Use()
				}
				optsInfo += fmt.Sprintf("%d) Использовать %s (%s)\n", i, usable.Name(), usable.Description())
				i++
			}
		} else if cons, ok := item.(I.Consumable); ok {
			if cons.UsesLeft() > 0 {
				opts[i] = func() {
					cons.Consume()
				}
				optsInfo += fmt.Sprintf("%d) Использовать %s [x%d] (%s)\n", i, cons.Name(), cons.UsesLeft(), cons.Description())
				i++
			}
		}
	})

	s.Printf(optsInfo)

	canceled := s.receiveWithAlert(ctx, time.Minute, func(umsg I.UserMessage, cancel func()) {
		opt, err := strconv.Atoi(umsg.Message())
		if umsg.User().ID() == p.ID() && err == nil {
			if action := opts[opt]; action != nil {
				action()
				cancel()
			}
		}
	}, 45*time.Second, "Ещё 15 секунд, чтобы сделать выбор!")

	if !canceled {
		s.Printf("%s решает пропустить ход!", p.FullName())
	}
}

func (s *Session) makeEnemyAction(ctx context.Context, e I.Enemy, es *enemySquad.EnemySquad) {
	s.Printf("Ходит %s.", e.Name())
	if s.makePauseFor(ctx, 3*time.Second) != nil {
		return
	}

	p := s.players.ChooseRandomAlivePreferBlocking()
	dmg := p.Damage(e.Attack())

	if p.Alive() {
		s.Printf("%s атакует %s и наносит %d урона. (HP:%d)", e.Name(), p.FullName(), dmg, p.Health())
		if s.makePauseFor(ctx, time.Second) != nil {
			return
		}
	} else {
		s.Printf("%s убивает %s!", e.Name(), p.FullName())
		if s.makePauseFor(ctx, 5*time.Second) != nil {
			return
		}
	}
}
