extends Control

var profile: Profile

func _ready():
	profile = await $ProfileAPI.create_anonymous()
	if profile == null:
		$TopLabel.text = "невозможно создать временный профиль"
	else:
		attempt_connection()

func _button_clicked():
	if profile == null:
		profile = await $ProfileAPI.create_anonymous()
	if profile != null:
		attempt_connection()

func attempt_connection():
	$TopLabel.text = "Ожидание других смельчаков"
	var data = await $NewCompaignAPI.make_anonymous_request(NewCompaignAPI.Mode.SINGLE, profile.id)
	$TopLabel.text = "Команда подобрана!"
	$ServerAnswer.text = JSON.stringify(data)
