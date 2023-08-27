package mazegen

import (
	"context"
	"math/rand"
	"server/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestRegularGenerator_Newby - проверяет подземелье для новичков
func TestRegularGenerator_Newby(t *testing.T) {
	const testingSeed = 7431118648

	gen := mockBaseGenerator()
	gen.DungeonsCompletedAvg = 0
	gen.GuildScoreAvg = 0
	gen.FromOneGuild = false

	cases := []struct {
		playersCount uint
		expected     [][]string
	}{
		{
			playersCount: 1,
			expected: [][]string{
				{"trap", "treasure", "trap", "rest", "boss"},
				{"trap", "enemy", "treasure", "rest", "boss"},
				{"enemy", "trap", "treasure", "enemy", "rest", "boss"},
			},
		},
		{
			playersCount: 2,
			expected: [][]string{
				{"trap", "treasure", "enemy", "rest", "boss"},
				{"trap", "treasure", "trap", "enemy", "rest", "boss"},
				{"enemy", "trap", "treasure", "enemy", "rest", "boss"},
			},
		},
		{
			playersCount: 3,
			expected: [][]string{
				{"trap", "enemy", "enemy", "treasure", "rest", "boss"},
				{"trap", "treasure", "enemy", "enemy", "rest", "boss"},
				{"enemy", "trap", "treasure", "enemy", "rest", "boss"},
				{"enemy", "treasure", "enemy", "trap", "trap", "rest", "boss"},
			},
		},
		{
			playersCount: 4,
			expected: [][]string{
				{"trap", "enemy", "enemy", "treasure", "rest", "boss"},
				{"trap", "treasure", "enemy", "trap", "rest", "boss"},
				{"enemy", "treasure", "trap", "enemy", "trap", "rest", "boss"},
				{"treasure", "trap", "enemy", "enemy", "trap", "rest", "boss"},
			},
		},
	}

	for _, c := range cases {
		gen.ParticipantsCount = c.playersCount
		dungeon, err := gen.Generate(rand.NewSource(testingSeed))

		assert.Nil(t, err)
		assertDungeonStruct(t, dungeon, c.expected)
	}
}

// TestRegularGenerator_Mature - проверяет подземелье для бывалых игроков
func TestRegularGenerator_Mature(t *testing.T) {
	const testingSeed = 7431118648

	gen := mockBaseGenerator()
	gen.DungeonsCompletedAvg = 4 // opt / 2
	gen.GuildScoreAvg = 150      // opt / 2
	gen.FromOneGuild = false

	cases := []struct {
		playersCount uint
		expected     [][]string
	}{
		{
			playersCount: 1,
			expected: [][]string{
				{"trap", "trap", "enemy", "treasure", "rest", "boss"},
				{"trap", "treasure", "enemy", "enemy", "rest", "boss"},
				{"treasure", "trap", "treasure", "trap", "rest", "boss"},
				{"trap", "treasure", "enemy", "enemy", "trap", "rest", "boss"},
			},
		},
		{
			playersCount: 2,
			expected: [][]string{
				{"trap", "enemy", "enemy", "treasure", "rest", "boss"},
				{"trap", "treasure", "enemy", "treasure", "rest", "boss"},
				{"enemy", "treasure", "enemy", "enemy", "trap", "rest", "boss"},
				{"treasure", "trap", "trap", "enemy", "enemy", "rest", "boss"},
			},
		},
		{
			playersCount: 3,
			expected: [][]string{
				{"trap", "enemy", "enemy", "treasure", "rest", "boss"},
				{"trap", "trap", "treasure", "trap", "enemy", "rest", "boss"},
				{"treasure", "trap", "trap", "enemy", "trap", "rest", "boss"},
				{"enemy", "treasure", "enemy", "trap", "enemy", "rest", "boss"},
			},
		},
		{
			playersCount: 4,
			expected: [][]string{
				{"trap", "enemy", "trap", "enemy", "treasure", "rest", "boss"},
				{"trap", "enemy", "enemy", "trap", "treasure", "rest", "boss"},
				{"treasure", "trap", "treasure", "enemy", "enemy", "rest", "boss"},
				{"enemy", "trap", "treasure", "enemy", "trap", "rest", "boss"},
			},
		},
	}

	for _, c := range cases {
		gen.ParticipantsCount = c.playersCount
		dungeon, err := gen.Generate(rand.NewSource(testingSeed))

		assert.Nil(t, err)
		assertDungeonStruct(t, dungeon, c.expected)
	}
}

// TestRegularGenerator_Expert - проверяет подземелье для настоящих экспертов
func TestRegularGenerator_Expert(t *testing.T) {
	const testingSeed = 7431118648

	gen := mockBaseGenerator()
	gen.DungeonsCompletedAvg = 12 // opt * 1.5
	gen.GuildScoreAvg = 450       // opt * 1.5
	gen.FromOneGuild = false

	cases := []struct {
		playersCount uint
		expected     [][]string
	}{
		{
			playersCount: 1,
			expected: [][]string{
				{"trap", "enemy", "trap", "enemy", "treasure", "rest", "boss"},
				{"trap", "enemy", "enemy", "enemy", "treasure", "rest", "boss"},
				{"treasure", "trap", "treasure", "enemy", "enemy", "rest", "boss"},
				{"enemy", "trap", "treasure", "enemy", "trap", "rest", "boss"},
				{"treasure", "enemy", "enemy", "enemy", "enemy", "trap", "rest", "boss"},
			},
		},
		{
			playersCount: 2,
			expected: [][]string{
				{"trap", "enemy", "enemy", "enemy", "treasure", "rest", "boss"},
				{"trap", "trap", "enemy", "trap", "treasure", "rest", "boss"},
				{"treasure", "trap", "enemy", "enemy", "enemy", "rest", "boss"},
				{"enemy", "enemy", "trap", "trap", "treasure", "trap", "rest", "boss"},
				{"enemy", "trap", "trap", "enemy", "enemy", "treasure", "rest", "boss"},
			},
		},
		{
			playersCount: 3,
			expected: [][]string{
				{"trap", "enemy", "treasure", "enemy", "treasure", "rest", "boss"},
				{"trap", "enemy", "enemy", "trap", "treasure", "rest", "boss"},
				{"enemy", "enemy", "trap", "trap", "enemy", "treasure", "rest", "boss"},
				{"trap", "enemy", "trap", "treasure", "trap", "enemy", "rest", "boss"},
				{"enemy", "trap", "treasure", "enemy", "enemy", "trap", "rest", "boss"},
				{"enemy", "enemy", "treasure", "trap", "trap", "enemy", "rest", "boss"},
			},
		},
		{
			playersCount: 4,
			expected: [][]string{
				{"trap", "enemy", "enemy", "enemy", "treasure", "rest", "boss"},
				{"trap", "enemy", "enemy", "treasure", "enemy", "trap", "rest", "boss"},
				{"enemy", "enemy", "trap", "treasure", "enemy", "enemy", "rest", "boss"},
				{"trap", "enemy", "enemy", "treasure", "trap", "enemy", "rest", "boss"},
				{"trap", "trap", "treasure", "enemy", "treasure", "trap", "rest", "boss"},
				{"enemy", "enemy", "treasure", "enemy", "trap", "enemy", "rest", "boss"},
			},
		},
	}

	for _, c := range cases {
		gen.ParticipantsCount = c.playersCount
		dungeon, err := gen.Generate(rand.NewSource(testingSeed))

		assert.Nil(t, err)
		assertDungeonStruct(t, dungeon, c.expected)
	}
}

func mockBaseGenerator() *RegularDungeonGenerator {
	return &RegularDungeonGenerator{
		treasureRoomGen: mockRoomGen("treasure"),
		trapRoomGen:     mockRoomGen("trap"),
		enemyRoomGen:    mockRoomGen("enemy"),
		npcRoomGen:      mockRoomGen("npc"),
		restRoomGen:     mockRoomGen("rest"),
		bossRoomGen:     mockRoomGen("boss"),
	}
}

func mockRoomGen(roomName string) RoomGenFunc {
	return func(s rand.Source) (domain.Room, error) {
		return &MockRoom{Name: roomName}, nil
	}
}

type MockRoom struct {
	Name string
}

func (*MockRoom) Visit(context.Context, *domain.Compaign) error {
	panic("must be never called")
}

func assertDungeonStruct(t *testing.T, dungeon domain.Dungeon, floors [][]string) {
	rd := dungeon.(*RegularDungeon)
	assert.Equal(t, len(floors), len(rd.floors))
	for i := range floors {
		assert.Equal(t, len(floors[i]), len(rd.floors[i]))
		for j := range floors[i] {
			assert.Equal(t, floors[i][j], rd.floors[i][j].(*MockRoom).Name)
		}
	}
}
