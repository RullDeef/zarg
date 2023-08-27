package domain

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestFight_Normal - обычный тест
func TestFight_Normal(t *testing.T) {
	actionChan := make(chan string)
	left := []Fightable{
		newMockFightable("p1", 1, actionChan),
		newMockFightable("p2", 1, actionChan),
		newMockFightable("p3", 1, actionChan),
	}
	right := []Fightable{
		newMockFightable("e1", 1, actionChan),
		newMockFightable("e2", 1, actionChan),
		newMockFightable("e3", 1, actionChan),
	}
	order := []Fightable{left[0], right[0], left[1], right[1], left[2], right[2]}

	fight, err := NewFight(left, right, order)
	assert.Nil(t, err)

	go func() {
		err := fight.PerformFight(context.Background())
		close(actionChan)
		assert.Nil(t, err)
	}()

	actions := make([]string, 0, 6)
	for action := range actionChan {
		actions = append(actions, action)
	}

	assert.Equal(t, []string{"p1", "e1", "p2", "e2", "p3"}, actions)
}

// TestFight_Chilling - некоторые враги не встречаются в очередности
func TestFight_Chilling(t *testing.T) {
	actionChan := make(chan string)
	left := []Fightable{
		newMockFightable("p1", 1, actionChan),
		newMockFightable("p2", 1, actionChan), // chilling
		newMockFightable("p3", 1, actionChan),
	}
	right := []Fightable{
		newMockFightable("e1", 1, actionChan), // chilling
		newMockFightable("e2", 1, actionChan),
		newMockFightable("e3", 1, actionChan),
	}
	order := []Fightable{left[0], right[1], left[2], right[2]}

	fight, err := NewFight(left, right, order)
	assert.Nil(t, err)

	go func() {
		err := fight.PerformFight(context.Background())
		close(actionChan)
		assert.Nil(t, err)
	}()

	actions := make([]string, 0, 4)
	for action := range actionChan {
		actions = append(actions, action)
	}

	assert.Equal(t, []string{"p1", "e2", "p3", "e3"}, actions)
}

// TestFight_MultipleMoves - некоторые враги встречаются в очередности несколько раз
func TestFight_MultipleMoves(t *testing.T) {
	actionChan := make(chan string)
	left := []Fightable{
		newMockFightable("p1", 1, actionChan),
		newMockFightable("p2", 3, actionChan),
		newMockFightable("p3", 2, actionChan),
	}
	right := []Fightable{
		newMockFightable("e1", 1, actionChan),
		newMockFightable("e2", 2, actionChan),
		newMockFightable("e3", 1, actionChan),
	}
	order := []Fightable{
		left[0], right[0],
		left[1], left[1], right[1],
		left[2], right[2], left[1], right[2],
	}

	fight, err := NewFight(left, right, order)
	assert.Nil(t, err)

	go func() {
		err := fight.PerformFight(context.Background())
		close(actionChan)
		assert.Nil(t, err)
	}()

	actions := make([]string, 0, 4)
	for action := range actionChan {
		actions = append(actions, action)
	}

	assert.Equal(t, []string{"p1", "e1", "p2", "p2", "e2", "p3", "e3", "p2", "e2"}, actions)
}

// newMockFightable - создает тестовую сущность, которая
// при каждом своем ходе отправляет свое имя в канал
func newMockFightable(name string, hp int, actions chan string) Fightable {
	return &mockFightable{
		name:    name,
		hp:      hp,
		actions: actions,
	}
}

type mockFightable struct {
	name    string
	hp      int
	actions chan string
}

func (m *mockFightable) IsAlive() bool {
	return m.hp > 0
}

func (m *mockFightable) MakeMove(_ context.Context) error {
	m.hp--
	m.actions <- m.name
	return nil
}
