class_name NewCompaignAPI
extends Node

const _websocketUrl = "ws://localhost:4667/compaigns/new"

enum Mode { SINGLE, RANDOM, GUILD }

var _client: WebSocketPeer = null
var _requestToSend = null

signal _response_came(PackedByteArray)

class Response:
	var compaign_id: String
	var players: Array

func make_anonymous_request(mode: Mode, profileID: String) -> Response:
	_client = WebSocketPeer.new()
	_requestToSend = {
		"mode": _mode_to_string(mode),
		"anonymous": true,
		"profile_id": profileID,
	}
	_client.connect_to_url(_websocketUrl)
	var data = await _response_came
	data = JSON.parse_string(data.get_string_from_utf8())
	var response = Response.new()
	response.compaign_id = data["compaign_id"]
	response.players = data["players"]
	return response

func _process(_delta):
	if _client == null:
		return
	_client.poll()
	var state = _client.get_ready_state()
	match state:
		WebSocketPeer.STATE_OPEN:
			# send request with provided mode and profile
			if _requestToSend != null:
				_client.send_text(JSON.stringify(_requestToSend))
				_requestToSend = null
			# wait until response comes
			if _client.get_available_packet_count() > 0:
				var packet = _client.get_packet()
				# emit signal
				_response_came.emit(packet)
				# then close client
				_client.close()
		WebSocketPeer.STATE_CLOSED:
			_client = null

func _mode_to_string(mode: Mode) -> String:
	match mode:
		Mode.SINGLE:
			return "single"
		Mode.RANDOM:
			return "random"
		Mode.GUILD:
			return "guild"
		_:
			push_error("unknown mode", mode)
			return ""
