package fight

import (
	"math/rand"
	I "zarg/lib/model/interfaces"
)

type Manager struct {
	entityAction func(e I.Entity, friends, oponents I.EntityList)
}

func New(entityAction func(e I.Entity, friends, oponents I.EntityList)) *Manager {
	return &Manager{
		entityAction,
	}
}

func (m *Manager) PerformFight(squad1, squad2 I.EntityList) int {
	if squad1.LenAlive() == 0 || squad2.LenAlive() == 0 {
		panic("expected to have someone alive in each squad")
	}

	// iterate over rounds
	for squad1.LenAlive() > 0 && squad2.LenAlive() > 0 {
		order := m.populateEntities(squad1, squad2)

		// iterate over turns
		for _, e := range order {
			friends := squad1
			oponents := squad2
			if squad2.Has(e) {
				friends = squad2
				oponents = squad1
			}

			// perform action
			m.entityAction(e, friends, oponents)

			// early end
			if squad1.LenAlive() == 0 || squad2.LenAlive() == 0 {
				break
			}
		}
	}

	if squad1.LenAlive() == 0 && squad2.LenAlive() == 0 {
		return 0
	} else if squad2.LenAlive() == 0 {
		return 1
	} else {
		return 2
	}
}

func (m *Manager) populateEntities(squad1, squad2 I.EntityList) []I.Entity {
	var order []I.Entity

	// populate entities
	squad1.ForEachAlive(func(e I.Entity) {
		order = append(order, e)
	})
	squad2.ForEachAlive(func(e I.Entity) {
		order = append(order, e)
	})
	rand.Shuffle(len(order), func(i, j int) {
		order[i], order[j] = order[j], order[i]
	})

	return order
}
