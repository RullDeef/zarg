package interfaces

var (
	StatusEffectStun = func(n int) StatusEffect {
		return StatusEffect{
			Name:        "🌀",
			Description: "оглушение. Пропуск хода",
			TimeLeft:    n,
		}
	}
	StatusEffectAgility = func(n int) StatusEffect {
		return StatusEffect{
			Name:        "⚡",
			Description: "проворность. Дополнительный ход",
			TimeLeft:    n,
		}
	}
	StatusEffectRegeneration = func(n int) StatusEffect {
		return StatusEffect{
			Name:        "💞",
			Description: "регенерация. +5❤",
			TimeLeft:    n,
		}
	}
	StatusEffectBleeding = func(n int) StatusEffect {
		return StatusEffect{
			Name:        "❣",
			Description: "кровотечение. -1❤",
			TimeLeft:    n,
		}
	}
	StatusEffectBurning = func(n int) StatusEffect {
		return StatusEffect{
			Name:        "🔥",
			Description: "горение. -2❤",
			TimeLeft:    n,
		}
	}
	StatusEffectFreezing = func(n int) StatusEffect {
		return StatusEffect{
			Name:        "❄",
			Description: "обморожение. -2❤",
			TimeLeft:    n,
		}
	}
	StatusEffectWeakness = func(n int) StatusEffect {
		return StatusEffect{
			Name:        "⚓",
			Description: "слабость. x0.8🗡",
			TimeLeft:    n,
		}
	}
	StatusEffectStrongness = func(n int) StatusEffect {
		return StatusEffect{
			Name:        "💪",
			Description: "сила. x1.25🗡",
			TimeLeft:    n,
		}
	}
)
