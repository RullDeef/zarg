class_name Profile
extends Object

var id: String
var nickname: String
var avatar: String # url
var money: int
var strength: int
var endurance: int
var luck: int
var observation: int
var inventory: Inventory

static func fromJSON(jsonStr) -> Profile:
	var profile = Profile.new()
	var parse_result = JSON.parse_string(jsonStr)
	match typeof(parse_result):
		TYPE_DICTIONARY:
			profile.id = parse_result["id"]
			profile.nickname = parse_result["nickname"]
			profile.avatar = parse_result["avatar"]
			profile.money = parse_result["money"]
			profile.strength = parse_result["strength"]
			profile.endurance = parse_result["endurance"]
			profile.luck = parse_result["luck"]
			profile.observation = parse_result["observation"]
			profile.inventory = Inventory.fromDict(parse_result["inventory"])
		_:
			push_error("unexpected json object type")
			return null
	return profile
