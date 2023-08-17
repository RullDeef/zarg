package domain

import "context"

// Entity - базовая сущность, представляющая игрока, врага (или союзника?)
type Entity struct {
	Name      string // имя сущности
	MaxHealth int    // максимальное здоровье
	Health    int    // текущее здоровье

	*EffectGroup // группа эффектов сущности
	*Inventory   // инвентарь сущности (враги тоже могут пользоваться предметами)

	// хуки для ключевых событий, которые могут возникнуть с существом
	onDeath          func() // срабатывает при смерти существа
	onBeforeHit      func() // срабатывает перед получением урона
	onBeforeDeathHit func() // срабатывает перед смертельным уроном
	onAfterHit       func() // срабатывает после получения несмертельного урона
	onAfterHeal      func() // срабатывает после восстановлении здоровья

	onMove func() error // ход существа, после срабатывания всех эффектов
}

// NewEntity - создает новую сущность в качестве базы для игрока, врага или союзника
func NewEntityBase(name string, health int, inventory *Inventory, effectGroup *EffectGroup) Entity {
	return Entity{
		Name:      name,
		MaxHealth: health,
		Health:    health,

		Inventory:   inventory,
		EffectGroup: effectGroup,

		onDeath:          func() {},
		onBeforeHit:      func() {},
		onBeforeDeathHit: func() {},
		onAfterHit:       func() {},
		onAfterHeal:      func() {},

		onMove: func() error { return nil },
	}
}

// IsAlive - возвращает true, если сущность жива
func (e *Entity) IsAlive() bool {
	return e.Health > 0
}

// Damage - наносит урон существу
func (e *Entity) Damage(damage int) {
	if damage < 0 {
		panic("negative damage")
	} else if e.Health == 0 {
		return
	} else if e.Health < damage {
		e.onBeforeDeathHit()
		e.Health = 0
		e.onDeath()
	} else {
		e.onBeforeHit()
		e.Health -= damage
		e.onAfterHit()
	}
}

// Heal - восстанавливает здоровье существу
func (e *Entity) Heal(heal int) {
	if e.Health += heal; e.Health > e.MaxHealth {
		e.Health = e.MaxHealth
	}
	e.onAfterHeal()
}

// MakeMove - реализация Fightable по-умолчанию (применяет эффекты и вызывает e.onMove)
func (e *Entity) MakeMove(context.Context) error {
	err := e.ActivateEffects()
	if err != nil {
		return err
	}
	return e.onMove()
}
