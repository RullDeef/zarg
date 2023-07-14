extends Control

var chatClient: TextChatClient = TextChatClient.new()

func _enter_tree():
	Logger.set_target(self)
	chatClient.chatID = CompaignHolder.textchatID
	chatClient.userID = CompaignHolder.profileID
	chatClient.chat_shutdown_signal.connect(_on_chat_shutdown)
	chatClient.new_message_signal.connect(_on_new_message)
	chatClient.user_connected_signal.connect(_on_user_connected)
	chatClient.user_disconnected_signal.connect(_on_user_disconnected)
	add_child(chatClient)
	Logger.log_info("chatID: %s" % [chatClient.chatID])
	$ChatIDLabel.text = "chatID: %s" % [chatClient.chatID]

func _exit_tree():
	chatClient.chat_shutdown_signal.disconnect(_on_chat_shutdown)
	chatClient.new_message_signal.disconnect(_on_new_message)
	chatClient.user_connected_signal.disconnect(_on_user_connected)
	chatClient.user_disconnected_signal.disconnect(_on_user_disconnected)
	Logger.detach_target()

func _on_send_button_pressed():
	var msg = $MsgEntry.text
	chatClient.send_message(msg)
	$MsgEntry.text = ""

func _on_chat_shutdown(reason: String):
	$RichTextLabel.add_text("[chat was shutdown: %s]\n" % [reason])

func _on_new_message(msg: TextChatClient.TextMessage):
	$RichTextLabel.add_text("> %s\n" % [msg.body])

func _on_user_connected(userID: String):
	$RichTextLabel.add_text("@newUser: %s\n" % [userID])

func _on_user_disconnected(userID: String):
	$RichTextLabel.add_text("@userLeave: %s\n" % [userID])

func _on_back_button_pressed():
	# go to compaign scene
	get_tree().change_scene_to_file("res://scenes/compaign/wait.tscn")
