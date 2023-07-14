class_name ProfileAPI
extends Node

const _serverUrlFmt = "http://%s:4667/profiles"

signal request_finished(PackedByteArray)

func create_anonymous() -> Profile:
	var _handler = func(_result, _response_code, _headers, body):
		request_finished.emit(body)
	var req = HTTPRequest.new()
	add_child(req)
	req.connect("request_completed", _handler)
	var err = req.request("%s/new" % [_serverUrlFmt % [GlobalConfig.serverHost]], PackedStringArray(), HTTPClient.METHOD_POST)
	if err != OK:
		push_error("failed to create request", req)
		remove_child(req)
		return null
	var body: PackedByteArray = await request_finished
	remove_child(req)
	return Profile.fromJSON(body.get_string_from_utf8())
