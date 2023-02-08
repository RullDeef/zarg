package floormaze

import (
	"log"
	"math/rand"
	"zarg/lib/model/enemy"
	"zarg/lib/model/enemy/boss"
	enemySquad "zarg/lib/model/enemy/squad"
	I "zarg/lib/model/interfaces"
	"zarg/lib/model/items/armor"
	"zarg/lib/model/items/potion"
	"zarg/lib/model/trap"
	"zarg/lib/model/weapon"
	"zarg/lib/utils"
)

func GenFloorMaze(name string) *FloorMaze {
	var rooms []any

	enemyRoomsCount := 3
	treasureRoomsCount := 2 + rand.Intn(2)

	for enemyRoomsCount > 0 || treasureRoomsCount > 0 {
		pm := utils.NewPropMap()
		pm.Add("enemy", enemyRoomsCount)
		pm.Add("treasure", treasureRoomsCount)

		switch pm.Choose().(string) {
		case "enemy":
			rooms = append(rooms, genEnemyRoom())
			enemyRoomsCount -= 1
		case "treasure":
			if rand.Float32() < 0.5 {
				rooms = append(rooms, genTrapRoom(), genTreasureRoom())
			} else {
				rooms = append(rooms, genEnemyRoom(), genTreasureRoom())
			}
			treasureRoomsCount -= 1
		}
	}

	rooms = append(rooms, genRestRoom(), genBossRoom())
	return newFloorMaze(name, rooms)
}

func genEnemyRoom() *EnemyRoom {
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

	return &EnemyRoom{
		Enemies: enemies,
	}
}

func genTrapRoom() *TrapRoom {
	probMap := utils.NewPropMap()

	probMap.Add(trap.New("Гигантская стрела вылетела прямо из стены!", trap.DamageRandom, 15), 4)
	probMap.Add(trap.New("Острые шипы выступают перед вашими ногами!", trap.DamageFirst, 22), 4)
	probMap.Add(trap.New("С потолка сваливается огромный камень!", trap.DamageRandom, 11), 5)
	probMap.Add(trap.New("Из темноты вылетает стая летучих мышей!", trap.DamageEveryone, 9), 5)
	probMap.Add(trap.New("Вы попадаете под душ из кислоты!", trap.DamageEveryone, 13), 2)
	probMap.Add(trap.New("Пол проваливается, и кто-то оказывается в лаве!", trap.DamageRandom, 999), 1)

	t := probMap.Choose().(*trap.Trap)

	return &TrapRoom{
		Trap: t,
	}
}

func genTreasureRoom() *TreasureRoom {
	probMap := utils.NewPropMap()
	for i := 0; i < 3; i++ {
		probMap.Add(weapon.RandomWeapon(10, 4), 1)
	}
	for i := 0; i < 3; i++ {
		probMap.Add(armor.Random(), 1)
	}
	for i := 0; i < 3; i++ {
		probMap.Add(potion.Random(), 1)
	}

	var items []I.Pickable
	for i := 0; i < 6; i++ {
		item := probMap.Choose()
		items = append(items, item.(I.Pickable))
		probMap.Add(item, -1)
	}

	return &TreasureRoom{
		Items: items,
	}
}

func genRestRoom() *RestRoom {
	return &RestRoom{}
}

func genBossRoom() *BossRoom {
	boss := boss.New(
		boss.NewPhase("Тролль", 300, func() I.DamageStats {
			return I.DamageStats{
				Base:       30,
				Crit:       60,
				CritChance: 0.25,
			}
		}, func(bp1, bp2 *boss.BossPhase) {
			log.Print("TODO: Троль разгневался и стал сильнее!")
		}),
		boss.NewPhase("Разъяренный Тролль", 200, func() I.DamageStats {
			return I.DamageStats{
				Base:       50,
				Crit:       80,
				CritChance: 0.4,
			}
		}, nil),
	)

	return &BossRoom{
		Boss: boss,
	}
}
