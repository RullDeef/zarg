package session

import (
	"context"
	"fmt"
	"strings"
	"time"
	"zarg/lib/model"
	I "zarg/lib/model/interfaces"
	"zarg/lib/model/trap"
	"zarg/lib/utils"
)

func (s *Session) exploreTrapRoom(ctx context.Context, fm *model.FloorMaze) {
	s.interactor.Printf("Что-то подсказывает вам, что тут не все так безобидно, как кажется на первый взгляд.")

	if s.makePauseFor(ctx, 4*time.Second) != nil {
		return
	}

	probMap := utils.NewPropMap()

	probMap.Add(trap.New("Гигантская стрела вылетела прямо из стены!", trap.DamageRandom, 15), 4)
	probMap.Add(trap.New("Острые шипы выступают перед вашими ногами!", trap.DamageFirst, 22), 4)
	probMap.Add(trap.New("С потолка сваливается огромный камень!", trap.DamageRandom, 11), 5)
	probMap.Add(trap.New("Из темноты вылетает стая летучих мышей!", trap.DamageEveryone, 9), 5)
	probMap.Add(trap.New("Вы попадаете под душ из кислоты!", trap.DamageEveryone, 13), 2)
	probMap.Add(trap.New("Пол проваливается, и кто-то оказывается в лаве!", trap.DamageRandom, 999), 1)

	t := probMap.Choose().(*trap.Trap)
	s.interactor.Printf(t.Name())

	if s.makePauseFor(ctx, 2*time.Second) != nil {
		return
	}

	healths := make(map[int]int)
	s.players.ForEachAlive(func(p I.Player) {
		healths[p.ID()] = p.Health()
	})
	damaged := t.Activate(s.players)
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
