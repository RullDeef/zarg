package floormaze

import (
	"log"
	"math"
	"math/rand"
	"zarg/lib/model/enemy"
	"zarg/lib/model/enemy/boss"
	enemySquad "zarg/lib/model/enemy/squad"
	I "zarg/lib/model/interfaces"
	"zarg/lib/model/items/armor"
	"zarg/lib/model/items/potion"
	"zarg/lib/model/items/weapon"
	"zarg/lib/model/trap"
	"zarg/lib/utils"
)

func GenFloorMaze(name string, balancer I.FloorGenBalancer) *FloorMaze {
	var rooms []any

	treasureRoomsCount := balancer.TreasureRoomsCount()
	enemyRoomsCount := balancer.EnemyRoomsCount()
	trapRoomsCount := balancer.TrapRoomsCount()

	for enemyRoomsCount > 0 || treasureRoomsCount > 0 || trapRoomsCount > 0 {
		pm := utils.NewPropMap()
		pm.Add("treasure", treasureRoomsCount)
		pm.Add("enemy", enemyRoomsCount)
		pm.Add("trap", trapRoomsCount)

		switch pm.Choose().(string) {
		case "treasure":
			if rand.Float32() < 0.5 && trapRoomsCount > 0 {
				rooms = append(rooms, genTrapRoom(balancer), genTreasureRoom(balancer))
				trapRoomsCount -= 1
			} else {
				rooms = append(rooms, genEnemyRoom(balancer), genTreasureRoom(balancer))
				enemyRoomsCount -= 1
			}
			treasureRoomsCount -= 1
		case "enemy":
			rooms = append(rooms, genEnemyRoom(balancer))
			enemyRoomsCount -= 1
		case "trap":
			rooms = append(rooms, genTrapRoom(balancer))
			trapRoomsCount -= 1
		}
	}

	rooms = append(rooms, genRestRoom(), genBossRoom(balancer))
	return newFloorMaze(name, rooms)
}

func genEnemyRoom(balancer I.FloorGenBalancer) *EnemyRoom {
	eb := balancer.EnemyBalancer()

	enemies := enemySquad.New(eb.EnemiesCount(), func() I.Enemy {
		healthMin, healthMax := eb.Health()
		health := healthMin + rand.Intn(healthMax-healthMin+1)
		attackMin, attackMax := eb.Attack()
		attack := attackMin + rand.Intn(attackMax-attackMin+1)
		crit := int(math.Ceil(float64(attack) * float64(eb.ExtraCrit())))
		critChance := eb.CritChance()

		return enemy.Random(health, func() I.DamageStats {
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

func genTrapRoom(balancer I.FloorGenBalancer) *TrapRoom {
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

func genTreasureRoom(balancer I.FloorGenBalancer) *TreasureRoom {
	totalItems := balancer.ItemsInTreasureRoomCount()
	halfTotal := totalItems/2 + totalItems%2

	items := make([]I.Pickable, 0, 3*halfTotal)
	for i := 0; i < halfTotal; i++ {
		items = append(items, weapon.RandomWeapon(10, 4))
	}
	for i := 0; i < halfTotal; i++ {
		items = append(items, armor.Random())
	}
	for i := 0; i < halfTotal; i++ {
		items = append(items, potion.Random())
	}

	rand.Shuffle(3*halfTotal, func(i, j int) {
		items[i], items[j] = items[j], items[i]
	})

	return &TreasureRoom{
		Items: items[:totalItems],
	}
}

func genRestRoom() *RestRoom {
	return &RestRoom{}
}

func genBossRoom(balancer I.FloorGenBalancer) *BossRoom {
	eb := balancer.EnemyBalancer()

	healthMin, healthMax := eb.Health()
	health := healthMin + rand.Intn(healthMax-healthMin+1)
	health *= 5
	health = 10 * int(math.Ceil(float64(health)/10.0)) // round up

	attackMin, attackMax := eb.Attack()
	attackPhase1 := 3*attackMin + 2*rand.Intn(attackMax-attackMin+1)
	attackPhase2 := attackPhase1 + 3*rand.Intn(attackMax-attackMin+1)

	boss := boss.New(
		boss.NewPhase("Тролль", int(0.6*float32(health)), func() I.DamageStats {
			return I.DamageStats{
				Base:       attackPhase1,
				Crit:       int(float32(attackPhase1) * 1.2),
				CritChance: 0.25,
			}
		}, func(bp1, bp2 *boss.BossPhase) {
			log.Print("TODO: Троль разгневался и стал сильнее!")
		}),
		boss.NewPhase("Разъяренный Тролль", int(0.4*float32(health)), func() I.DamageStats {
			return I.DamageStats{
				Base:       attackPhase2,
				Crit:       int(float32(attackPhase2) * 1.2),
				CritChance: 0.4,
			}
		}, nil),
	)

	return &BossRoom{
		Boss: boss,
	}
}
