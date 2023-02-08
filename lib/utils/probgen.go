package utils

import (
	"math/rand"
	"strconv"
)

type ProbMap map[any]int // probability map

func NewPropMap() *ProbMap {
	return &ProbMap{}
}

func (p *ProbMap) Add(key any, val int) {
	if val != 0 {
		(*p)[key] += val
	}
}

func (p *ProbMap) Choose() any {
	sum := 0
	for _, val := range *p {
		sum += val
	}

	target := rand.Intn(sum + 1)
	for key, val := range *p {
		target -= val
		if target <= 0 {
			return key
		}
	}

	panic("must never occur " + strconv.Itoa(target))
}
