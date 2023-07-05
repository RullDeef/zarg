class_name Inventory
extends Object

var items: Array

func _init(_items: Array):
	items = _items

static func fromJSON(jsonStr: String) -> Inventory:
	var _items = JSON.parse_string(jsonStr)
	return Inventory.new(_items)

static func fromDict(dict: Dictionary) -> Inventory:
	return Inventory.new(dict["items"])
