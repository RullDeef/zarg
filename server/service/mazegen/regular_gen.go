package mazegen

import (
	"errors"
	"math"
	"math/rand"
	"server/domain"
)

// Константы, используемые при генерации подземелья
const (
	// Генерация количества этажей
	fMin    = 3    // F_Min
	fMax    = 8    // F_Max
	kFactor = 0.15 // k
	nOpt    = 4    // N_Opt
	rOpt    = 300  // R_Opt
	cOpt    = 8    // C_Opt

	// Генерация количества комнат на этаже
	rMin = 4 // R_Min
	rMax = 8 // R_Max
)

var (
	// ErrTooFewRoomsOnFloor - слишком мало комнат на этаже
	ErrTooFewRoomsOnFloor = errors.New("too few rooms on floor")

	// ErrAllRoomsInitialized - невозможно выбрать неинициализированную комнату
	ErrAllRoomsInitialized = errors.New("all rooms initialized already")

	// ErrInvalidNumWeight - некорректный элемент массива для взвешенного выбора
	ErrInvalidNumWeight = errors.New("invalid num weight")

	// ErrAllWeightsAreZero - сумма всех весов оказалась равна нулю
	ErrAllWeightsAreZero = errors.New("all weights are zero")
)

// RegularDungeonGenerator - генератор линейного (обычного) подземелья из нескольких этажей
type RegularDungeonGenerator struct {
	ParticipantsCount    uint    // N - количество человек в походе
	GuildScoreAvg        float64 // R - очки прогресса гильдии (если игроки из разных гильдий - среднее значение)
	DungeonsCompletedAvg float64 // C - среднее количество пройденных подземелий

	// FromOneGuild - являются ли игроки членами одной гильдии.
	// Данный параметр влияет на помещение НИПа в подземелье
	FromOneGuild bool

	treasureRoomGen RoomGenFunc // генератор комнат сокровищ
	trapRoomGen     RoomGenFunc // генератор комнат с ловушками
	enemyRoomGen    RoomGenFunc // генератор комнат с врагами
	restRoomGen     RoomGenFunc // генератор комнат отдыха
	bossRoomGen     RoomGenFunc // генератор комнат с боссом
	npcRoomGen      RoomGenFunc // генератор комнаты с НИПом
}

// RoomGenFunc - функция генерации комнаты определенного типа
type RoomGenFunc func(rand.Source) (domain.Room, error)

func (g *RegularDungeonGenerator) Generate(src rand.Source) (domain.Dungeon, error) {
	floors := make([][]domain.Room, g.floorsCount())

	for i := range floors {
		if floor, err := g.generateFloor(src, i); err != nil {
			return nil, err
		} else {
			floors[i] = floor
		}
	}

	// поместить НИПа
	if g.FromOneGuild && g.npcTest(src) {
		var err error
		floors, err = g.placeNPC(src, floors)
		if err != nil {
			return nil, err
		}
	}

	return &RegularDungeon{floors: floors}, nil
}

func (g *RegularDungeonGenerator) floorsCount() int {
	x := float64(g.ParticipantsCount) / nOpt
	x += g.GuildScoreAvg / rOpt
	x += g.DungeonsCompletedAvg / cOpt

	x = math.Tanh(kFactor * x)
	x = fMin + (fMax-fMin)*x // lerp

	return int(math.Round(x))
}

func (g *RegularDungeonGenerator) roomsCount(floor int) int {
	x := 1 + float64(floor)/2
	x += 2 * float64(g.ParticipantsCount) / nOpt
	x += g.GuildScoreAvg / rOpt
	x += g.DungeonsCompletedAvg / cOpt

	x = math.Log2(x) / 3
	x = rMin + (rMax-rMin)*x // lerp

	return int(math.Round(x))
}

func (g *RegularDungeonGenerator) generateFloor(src rand.Source, floor int) ([]domain.Room, error) {
	roomsCount := g.roomsCount(floor)
	if roomsCount < 4 {
		return nil, ErrTooFewRoomsOnFloor
	}

	rooms := make([]domain.Room, roomsCount)

	// последняя комната - комната с боссом
	if room, err := g.bossRoomGen(src); err != nil {
		return nil, err
	} else {
		rooms[roomsCount-1] = room
	}

	// предпоследняя комната - комната отдыха
	if room, err := g.restRoomGen(src); err != nil {
		return nil, err
	} else {
		rooms[roomsCount-2] = room
	}

	// случайно поместить сокровищницу
	if i, err := pickRandNil(src, rooms); err != nil {
		return nil, err
	} else if room, err := g.treasureRoomGen(src); err != nil {
		return nil, err
	} else {
		rooms[i] = room
	}

	// случайно поместить комнату с ловушками
	if i, err := pickRandNil(src, rooms); err != nil {
		return nil, err
	} else if room, err := g.trapRoomGen(src); err != nil {
		return nil, err
	} else {
		rooms[i] = room
	}

	// инициализировать остальные комнаты
	Rtreasure := 1
	Rtrap := 1
	Renemy := 0
	for {
		i, err := pickRandNil(src, rooms)
		if err == ErrAllRoomsInitialized {
			break
		}

		X := max(0, min(10+int(g.ParticipantsCount)-Rtreasure, 14))
		Y := max(0, min(2+int(g.DungeonsCompletedAvg)-Rtrap, 8))
		Z := max(0, min(30+floor-Renemy, 40))

		rType, err := pickWeighted(src, []int{X, Y, Z})
		if err != nil {
			return nil, err
		}

		var room domain.Room
		switch rType {
		case 0: // сокровищница
			room, err = g.treasureRoomGen(src)
			Rtreasure++
		case 1: // комната с ловушками
			room, err = g.trapRoomGen(src)
			Rtrap++
		case 2: // комната с монстрами
			room, err = g.enemyRoomGen(src)
			Renemy++
		default:
			panic("must never happen")
		}

		if err != nil {
			return nil, err
		}
		rooms[i] = room
	}

	return rooms, nil
}

// placeNPC - Размещает НИПа случайным образом в подземелье
func (g *RegularDungeonGenerator) placeNPC(src rand.Source, floors [][]domain.Room) ([][]domain.Room, error) {
	floorIndex, err := g.npcFloor(src, len(floors))
	if err != nil {
		return nil, err
	}

	// выбираем случайную комнату, кроме комнаты отдыха и комнаты с боссом.
	// Предполагается, что данные комнаты будут последними на любом этаже
	roomIndex := rand.New(src).Intn(len(floors[floorIndex]) - 2)

	if room, err := g.npcRoomGen(src); err != nil {
		return nil, err
	} else {
		floors[floorIndex][roomIndex] = room
	}

	return floors, nil
}

// npcTest - выполняет тест и с вероятностью P_NPC возвращает true
func (g *RegularDungeonGenerator) npcTest(src rand.Source) bool {
	x := g.DungeonsCompletedAvg/cOpt - g.GuildScoreAvg/rOpt
	x = math.Log(max(0, x) / 3)
	prob := 0.5 + 0.5*math.Tanh(x)
	return rand.New(src).Float64() < prob
}

// npcFloor - случайно выбирает индекс этажа для размещения НИПа
func (g *RegularDungeonGenerator) npcFloor(src rand.Source, floorCount int) (int, error) {
	weights := make([]int, floorCount)
	for i := range weights {
		weights[i] = i
	}
	return pickWeighted(src, weights)
}

// pickRandNil - случайно выбирает неинициализированную комнату и возвращает ее индекс
func pickRandNil(src rand.Source, rooms []domain.Room) (int, error) {
	allNotNil := true
	for _, room := range rooms {
		if room == nil {
			allNotNil = false
			break
		}
	}
	if allNotNil {
		return 0, ErrAllRoomsInitialized
	}
	rnd := rand.New(src)
	for {
		if i := rnd.Intn(len(rooms)); rooms[i] == nil {
			return i, nil
		}
	}
}

// pickWeighted - случайно выбирает один элемент из списка в соответствии с его весом.
// Если в массиве присутствуют нулевые элементы, они никогда не будут выбраны
func pickWeighted(src rand.Source, nums []int) (int, error) {
	sum := 0
	for _, num := range nums {
		if num < 0 {
			return 0, ErrInvalidNumWeight
		}
		sum += num
	}
	if sum == 0 {
		return 0, ErrAllWeightsAreZero
	}
	rnd := rand.New(src).Intn(sum)
	i := 0
	for i < len(nums)-1 && (rnd > 0 || nums[i] == 0) {
		rnd -= nums[i]
		i++
	}
	return i, nil
}
