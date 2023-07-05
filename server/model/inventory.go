package model

import (
	"bytes"
	"container/list"
	"encoding/json"
)

type Inventory struct {
	Items *list.List
}

type QuickAccessAssignment struct {
	Slot   int
	ItemID string
}

func NewEmptyInventory() *Inventory {
	return &Inventory{
		Items: list.New(),
	}
}

func (inv Inventory) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("{\"items\":")

	itemsBytes, _ := marshalList(inv.Items)
	buffer.Write(itemsBytes)

	buffer.WriteByte('}')
	return buffer.Bytes(), nil
}

func marshalList(list *list.List) ([]byte, error) {
	buffer := bytes.NewBufferString("[")

	for el := list.Front(); el != nil; el = el.Next() {
		marshalled, err := json.Marshal(el.Value)

		if err != nil {
			return nil, err
		}

		buffer.WriteString(string(marshalled))

		if el.Next() != nil { // Has more elements, lets append a comma
			buffer.WriteRune(',')
		}
	}

	buffer.WriteString("]")
	return buffer.Bytes(), nil
}
