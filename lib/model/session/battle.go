package session

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"
	"zarg/lib/model"
	"zarg/lib/model/enemy"
	enemySquad "zarg/lib/model/enemy/squad"
	I "zarg/lib/model/interfaces"
	"zarg/lib/utils"
)

func (s *Session) exploreEnemiesRoom(ctx context.Context, fm *model.FloorMaze) {
	s.interactor.Printf("Вы не одни... На вас напали!")

	enemies := enemySquad.New(2+rand.Intn(2), func() I.Enemy {
		attackMin := 8
		attackMax := 14
		attack := attackMin + rand.Intn(attackMax-attackMin+1)
		crit := attack + 10
		critChance := 0.05 + 0.05*rand.Float32()

		return enemy.Random(func() I.DamageStats {
			return I.DamageStats{
				Base:       attack,
				Crit:       crit,
				CritChance: critChance,
			}
		})
	})

	s.PerformBattle(ctx, enemies)
}

func (s *Session) PerformBattle(ctx context.Context, es *enemySquad.EnemySquad) {

	// show battle overall info
	battleInfo := fmt.Sprintf("Игроки:\n%sВраги:\n%s", s.players.CompactInfo(), es.CompactInfo())
	s.interactor.Printf(battleInfo)

	for es.LenAlive() > 0 && s.players.LenAlive() > 0 {
		if s.makePauseFor(ctx, time.Second) != nil {
			return
		}

		// make turn generator
		turnGen := utils.NewPropMap()

		turnGen.Add("players", s.players.LenAlive())
		turnGen.Add("enemies", es.LenAlive())

		switch turnGen.Choose().(string) {
		case "players":
			// show battle overall info
			battleInfo := fmt.Sprintf("Игроки:\n%sВраги:\n%s", s.players.CompactInfo(), es.CompactInfo())
			s.interactor.Printf(battleInfo)
			p := s.players.ChooseNext()
			s.makePlayerAction(ctx, p, es)
		case "enemies":
			e := es.ChooseNext()
			s.makeEnemyAction(ctx, e, es)
		default:
			log.Fatal("must never happen!")
		}
	}

	if es.LenAlive() == 0 {
		s.interactor.Printf("Битва завершена. Все враги повержены!")
	} else {
		s.interactor.Printf("Битва завершена. Все игроки мертвы!")
	}
}

func (s *Session) makePlayerAction(ctx context.Context, p I.Player, es *enemySquad.EnemySquad) {
	s.interactor.Printf("Ходит %s.", p.FullName())

	if s.makePauseFor(ctx, 5*time.Second) != nil {
		return
	}

	optsInfo := "Варианты действия:\n"
	opts := map[int]func(){}
	i := 1
	es.ForEachAlive(func(e I.Enemy) {
		optsInfo += fmt.Sprintf("%d) Атаковать %s (HP: %d)\n", i, e.Name(), e.Health())
		opts[i] = func() {
			dmg := e.Damage(p.Weapon().Attack())
			if e.Alive() {
				s.interactor.Printf("%s атакует %s и наносит %d урона.", p.FullName(), e.Name(), dmg)
			} else {
				s.interactor.Printf("%s убивает %s", p.FullName(), e.Name())
			}
		}
		i += 1
	})
	s.interactor.Printf(optsInfo)

	canceled := s.receiveWithAlert(ctx, time.Minute, func(umsg model.UserMessage, cancel func()) {
		opt, err := strconv.Atoi(umsg.Message)
		if umsg.User.ID() == p.ID() && err == nil {
			if action := opts[opt]; action != nil {
				action()
				cancel()
			}
		}
	}, 45*time.Second, "Ещё 15 секунд, чтобы сделать выбор!")

	if !canceled {
		s.interactor.Printf("%s решает пропустить ход!", p.FullName())
	}
}

func (s *Session) makeEnemyAction(ctx context.Context, e I.Enemy, es *enemySquad.EnemySquad) {
	s.interactor.Printf("Ходит %s.", e.Name())
	if s.makePauseFor(ctx, 3*time.Second) != nil {
		return
	}

	p := s.players.ChooseRandomAlive()
	dmg := p.Damage(e.Attack())

	if p.Alive() {
		s.interactor.Printf("%s атакует %s и наносит %d урона. (HP:%d)", e.Name(), p.FullName(), dmg, p.Health())
		if s.makePauseFor(ctx, time.Second) != nil {
			return
		}
	} else {
		s.interactor.Printf("%s убивает %s!", e.Name(), p.FullName())
		if s.makePauseFor(ctx, 5*time.Second) != nil {
			return
		}
	}
}
