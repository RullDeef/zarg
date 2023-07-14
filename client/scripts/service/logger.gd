# class_name Logger
extends Node

# where to display logs
var _target: Control = null
var _rich_text: RichTextLabel = null
var _hide_button: Button = null

enum LogLevel { INFO, DEBUG, WARN, ERROR }

func _ready():
	get_viewport().size_changed.connect(_on_viewport_resized)

func _on_viewport_resized():
	if _target != null:
		_target.size = get_viewport().get_visible_rect().size

func set_target(new_target: Control):
	if _target != null:
		detach_target()
	_target = Control.new()
	_rich_text = RichTextLabel.new()
	# _rich_text.size_flags_horizontal = Control.SIZE_FILL
	# _rich_text.size_flags_vertical = Control.SIZE_EXPAND
	_rich_text.fit_content = true
	_hide_button = Button.new()
	_hide_button.button_down.connect(_on_hide_button_press)
	_target.add_child(_rich_text)
	_target.add_child(_hide_button)
	new_target.add_child(_target)
	_target.set_anchors_preset(Control.PRESET_CENTER_TOP)
	_rich_text.set_anchors_preset(Control.PRESET_FULL_RECT)
	_hide_button.set_anchors_preset(Control.PRESET_TOP_RIGHT)
	_hide_button.set_offsets_preset(Control.PRESET_TOP_RIGHT, Control.PRESET_MODE_KEEP_SIZE, 20)
	_hide_button.text = '[?]'
	_hide_button.size = Vector2(40, 40)
	_target.get_parent().move_child(_target, 0)
	_on_viewport_resized()

func _on_hide_button_press():
	_rich_text.visible = !_rich_text.visible

func detach_target():
	if _target != null:
		_target.get_parent().remove_child(_target)
		_hide_button.button_down.disconnect(_on_hide_button_press)
		_target = null
		_rich_text = null
		_hide_button = null
	else:
		push_error("_target is null already")

func log_info(msg: String):
	_log_message(LogLevel.INFO, Color.CYAN, msg)

func log_debug(msg: String):
	_log_message(LogLevel.DEBUG, Color.WHITE, msg)

func log_warn(msg: String):
	_log_message(LogLevel.WARN, Color.ORANGE, msg)

func log_error(msg: String):
	_log_message(LogLevel.ERROR, Color.RED, msg)

func _log_message(level: LogLevel, color: Color, msg: String):
	if _rich_text != null:
		_rich_text.push_color(color)
		_rich_text.add_text("[%s] %s" % [_level_to_string(level), msg])
		_rich_text.pop()

func _level_to_string(level: LogLevel) -> String:
	match level:
		LogLevel.INFO:
			return "INFO"
		LogLevel.DEBUG:
			return "DEBUG"
		LogLevel.WARN:
			return "WARN"
		LogLevel.ERROR:
			return "ERROR"
		_:
			push_error("unknown log level value: %s" % [level])
			return "????"
