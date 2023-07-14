extends Control

var profile: Profile

func _enter_tree():
	Logger.set_target(self)
	Logger.log_info("initialized")

func _exit_tree():
	Logger.detach_target()

func _button_clicked():
	if profile == null:
		profile = await $ProfileAPI.create_anonymous()
		if profile == null:
			$TopLabel.text = "невозможно создать временный профиль"
	if profile != null and not $NewCompaignAPI.is_making_request():
		attempt_connection()

func attempt_connection():
	$TopLabel.text = "Ожидание других смельчаков"
	$ServerAnswer.text = "waiting..."
	var data = await $NewCompaignAPI.make_anonymous_request(NewCompaignAPI.Mode.RANDOM, profile.id)
	if data != null:
		$TopLabel.text = "Команда подобрана!"
		$ServerAnswer.text = "compaign_id: " + data.compaign_id + "\n" + \
			("players count: %d" % [data.players.size()])
		CompaignHolder.profileID = profile.id
		CompaignHolder.compaignID = data.compaign_id
		CompaignHolder.textchatID = data.textchat_id
		# switch to chat screen
		await get_tree().create_timer(1.0).timeout
		get_tree().change_scene_to_file("res://scenes/compaign/chat.tscn")

func _on_server_addr_text_changed(new_text):
	GlobalConfig.serverHost = new_text
