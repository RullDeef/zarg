class_name Compaign
extends Object

var id: String
var players: Array
var dungeon: Dungeon

static func fromJSON(jsonStr: String) -> Compaign:
	var dict = JSON.parse_string(jsonStr)
	match typeof(dict):
		TYPE_DICTIONARY:
			var comp = Compaign.new()
			comp.id = dict["id"]
			comp.players = dict["players"]
			comp.dungeon = dict["dungeon"]
			return comp
		_:
			push_error("unexpected json type for dict", dict)
			return null
