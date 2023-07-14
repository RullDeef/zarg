class_name TextChatClient
extends Node

const _wsServerUrlFormat = "ws://%s:4668/textchats/%s?user_id=%s"

var chatID: String = ""
var userID: String = ""

var _client: WebSocketPeer = null
var _retry_attempts: int = 1

signal chat_shutdown_signal(reason: String)
signal user_connected_signal(userID: String)
signal user_disconnected_signal(userID: String)
signal new_message_signal(TextMessage)

class TextMessage:
	var from: String # profile id for now
	var body: String
	var time: Dictionary

func _enter_tree():
	if chatID == "" or userID == "":
		push_error("TextChatClient must be initialized before enter tree")
		return
	_reconnect()

func _exit_tree():
	if _client != null:
		_client.close()
		_client = null

func _process(_delta):
	if _client == null:
		return
	_client.poll()
	var state = _client.get_ready_state()
	match state:
		WebSocketPeer.STATE_OPEN:
			_process_messages()
			pass
		WebSocketPeer.STATE_CLOSED:
			Logger.log_info("peer closed")
			_client = null
			_reconnect()

func _reconnect():
	if _retry_attempts == 5:
		push_error("max retry count")
		Logger.log_error("max retry count limit")
		return
	_client = WebSocketPeer.new()
	var url = _wsServerUrlFormat % [GlobalConfig.serverHost, chatID, userID]
	print("connecting to url:", url)
	Logger.log_info("connecting to url: %s" % [url])
	var err = _client.connect_to_url(url)
	if err != OK:
		Logger.log_error("failed to connect: %s (attempt: %s)" % [err, _retry_attempts])
		_reconnect()
	else:
		_retry_attempts = 1

func send_message(msg: String):
	if _client == null:
		push_error("failed to send message: TextChatClient not connected")
		Logger.log_error("failed to send message: TextChatClient not connected")
		return
	var err = _client.send_text(msg)
	if err != OK:
		msg = "failed to send message: %s" % [err]
		push_error(msg)
		Logger.log_error(msg)

func _process_messages():
	while _client.get_available_packet_count() > 0:
		var packet = _client.get_packet()
		var data = JSON.parse_string(packet.get_string_from_utf8())
		_parse_event(data)

# internal usage
enum _EventType {
	TEXTCHAT_SHUTDOWN,
	MESSAGE_NEW,
	USER_CONNECTED,
	USER_DISCONNECTED
}

func _parse_event(eventData: Dictionary):
	Logger.log_info("got new chat event: %s" % [eventData])
	match int(eventData["type"]):
		_EventType.TEXTCHAT_SHUTDOWN:
			chat_shutdown_signal.emit(eventData["reason"])
		_EventType.USER_CONNECTED:
			user_connected_signal.emit(eventData["user_id"])
		_EventType.USER_DISCONNECTED:
			user_disconnected_signal.emit(eventData["user_id"])
		_EventType.MESSAGE_NEW:
			var msg = TextMessage.new()
			msg.time = Time.get_datetime_dict_from_datetime_string(eventData["time"], false)
			msg.from = eventData["from_id"]
			msg.body = eventData["body"]
			new_message_signal.emit(msg)
		_:
			Logger.log_error("unknown event: %s" % [eventData])
			push_error("unknown event: ", eventData)
