package model

type Item interface {
	ID() string
	Title() string
	Description() string

	Cost() int
	Cellable() bool

	Weight() float64
	Amount() int
}
