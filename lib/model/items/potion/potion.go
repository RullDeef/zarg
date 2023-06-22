package potion

import (
	I "zarg/lib/model/interfaces"
	"zarg/lib/utils"
)

func Random() I.Pickable {
	pm := utils.NewPropMap()

	pm.Add(NewHealingPotion("Зелье восстановления I", 20, 1), 3)
	pm.Add(NewHealingPotion("Зелье восстановления I", 20, 3), 2)
	pm.Add(NewHealingPotion("Зелье восстановления I", 20, 5), 1)
	pm.Add(NewHealingPotion("Зелье восстановления II", 35, 1), 2)
	pm.Add(NewHealingPotion("Зелье восстановления II", 35, 2), 1)

	pm.Add(NewStrengthPotion("Зелье силы I", 3, 1), 3)
	pm.Add(NewStrengthPotion("Зелье силы I", 3, 3), 2)
	pm.Add(NewStrengthPotion("Зелье силы I", 3, 5), 1)
	pm.Add(NewStrengthPotion("Зелье силы II", 5, 1), 2)
	pm.Add(NewStrengthPotion("Зелье силы II", 5, 2), 1)

	return pm.Choose().(I.Pickable)
}
