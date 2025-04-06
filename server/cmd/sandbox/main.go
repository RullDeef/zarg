// Sandbox - модуль песочницы для тестирования локальных гипотез и отловли багов.
// Код в песочнице не должен следовать принципам чистой архитектуры,
// его задача - тестирование гипотез.

package main

import (
	"server/domain"
)

func main() {
	// testRegularDungeon()
	testImporter()
}

// создает тестового игрока с пустым инвентарем
func newTestProfile(name string) *domain.Profile {
	profile := domain.NewAnonymousProfile()
	profile.Nickname = name
	return profile
}
