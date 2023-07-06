extends Control

var profile: Profile

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


func _on_server_addr_text_changed(new_text):
	$NewCompaignAPI.serverHost = new_text
	$ProfileAPI.serverHost = new_text
