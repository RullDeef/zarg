package domain

import (
	"errors"
	"fmt"
)

var ErrSameEffect = errors.New("same effect") // данный эффект уже присутствует в группе

// EffectGroup - контейнер для активных эффектов некоторого существа
type EffectGroup struct {
	// Effects - активные эффекты (только для чтения!)
	Effects []*Effect

	// policy - политика взаимодействия эффектов.
	// Используется при модификации списка активных эффектов при добавлении нового эффекта
	policy EffectPolicy
}

// EffectPolicy - политика взаимодействия эффектов.
// Некоторые эффекты не могут существовать вместе; данный интерфейс
// разрешает спорные моменты, определяя взаимодействия между
// имеющимися эффектами и новым накладываемым эффектом
type EffectPolicy interface {
	// Interaction - определяет результат взаимодействия эффектов
	Interaction(new *Effect, existing *Effect) (EffectInteractionResult, error)
}

type EffectInteractionResult int

const (
	EffectsCoexist        EffectInteractionResult = iota // эффекты разного типа, могут быть активны одновременно
	EffectsMerge                                         // эффекты одного типа, должны быть объеденены
	EffectsOpposite                                      // эффекты противоположных типов (взаимоуничтожаются)
	EffectsNewReplacesOld                                // новый эффект заменяет старый
	EffectsOldReplacesNew                                // старый эффект заменяет новый
)

type EffectPolicyFunc func(*Effect, *Effect) (EffectInteractionResult, error)

func (f EffectPolicyFunc) Interaction(e1 *Effect, e2 *Effect) (EffectInteractionResult, error) {
	return f(e1, e2)
}

// NewEffectGroup - создает новую группу эффектов с указанной политикой взаимодействия эффектов
func NewEffectGroup(policy EffectPolicy) *EffectGroup {
	return &EffectGroup{policy: policy}
}

// AddEffect - добавляет новый эффект в группу
func (g *EffectGroup) AddEffect(e *Effect) error {
	if g.hasEffect(e) {
		return ErrSameEffect
	}

	// check policies
	for i, effect := range g.Effects {
		res, err := g.policy.Interaction(e, effect)
		if err != nil {
			return err
		}

		switch res {
		case EffectsCoexist:
			continue
		case EffectsMerge:
			if newEf, err := g.Effects[i].MergeWith(e); err != nil {
				return err
			} else {
				g.Effects[i] = newEf
				return nil
			}
		case EffectsOpposite:
			g.Effects = append(g.Effects[:i], g.Effects[i+1:]...)
			return nil
		case EffectsNewReplacesOld:
			g.Effects[i] = e
			return nil
		case EffectsOldReplacesNew:
			return nil
		default:
			// must never happen
			panic(fmt.Errorf("unknown effect interaction result: %v", res))
		}
	}

	g.Effects = append(g.Effects, e)
	return nil
}

func (g *EffectGroup) hasEffect(e *Effect) bool {
	for _, effect := range g.Effects {
		if effect == e {
			return true
		}
	}
	return false
}

// ActivateEffects - активирует все эффекты в порядке их добавления.
// Также удаляет истекшие эффекты
func (g *EffectGroup) ActivateEffects() error {
	j := 0
	for _, effect := range g.Effects {
		expired, err := effect.Activate()
		if err != nil {
			return err
		}
		if !expired {
			g.Effects[j] = effect
			j++
		}
	}
	for i := j; i < len(g.Effects); i++ {
		g.Effects[i] = nil
	}
	g.Effects = g.Effects[:j]
	return nil
}

// NoopEffectPolicy - прозрачная политика (ничего не делает, все эффекты сохраняются)
func NoopEffectPolicy() EffectPolicy {
	return EffectPolicyFunc(func(e1 *Effect, e2 *Effect) (EffectInteractionResult, error) {
		return EffectsCoexist, nil
	})
}
