package textchat

import (
	"errors"
	"time"
)

type EventType int64

const (
	EventChatShutdownType = EventType(iota)
	EventMessageNewType
	EventUserConnectedType
	EventUserDisconnectedType
)

var (
	ErrFailedToParseEvent = errors.New("failed to parse event")
	ErrInvalidFieldType   = errors.New("invalid type for event field")
	ErrInvalidEventTime   = errors.New("unable to parse event time")
	ErrUnknownEventType   = errors.New("unknown event type")
)

type Event interface {
	Type() EventType
	Time() time.Time
}

type BaseEvent struct {
	EventType EventType `json:"type"`
	EventTime time.Time `json:"time"`
}

func (e *BaseEvent) Type() EventType {
	return e.EventType
}

func (e *BaseEvent) Time() time.Time {
	return e.EventTime
}

type EventChatShutdown struct {
	BaseEvent

	Reason string `json:"reason"`
}

func newEventChatShutdown(reason string) Event {
	return &EventChatShutdown{
		BaseEvent: BaseEvent{
			EventType: EventChatShutdownType,
			EventTime: time.Now(),
		},
		Reason: reason,
	}
}

type EventMessageNew struct {
	BaseEvent

	FromID string `json:"from_id"`
	Body   string `json:"body"`
}

func newEventMessageNew(fromID string, body string) Event {
	return &EventMessageNew{
		BaseEvent: BaseEvent{
			EventType: EventMessageNewType,
			EventTime: time.Now(),
		},
		FromID: fromID,
		Body:   body,
	}
}

type EventUserConnected struct {
	BaseEvent

	UserID string `json:"user_id"`
}

func newEventUserConnected(userID string) Event {
	return &EventUserConnected{
		BaseEvent: BaseEvent{
			EventType: EventUserConnectedType,
			EventTime: time.Now(),
		},
		UserID: userID,
	}
}

type EventUserDisconnected struct {
	BaseEvent

	UserID string `json:"user_id"`
}

func NewEventUserDisconnected(userID string) Event {
	return &EventUserDisconnected{
		BaseEvent: BaseEvent{
			EventType: EventUserDisconnectedType,
			EventTime: time.Now(),
		},
		UserID: userID,
	}
}

func ParseEvent(json map[string]any) (Event, error) {
	Type, err := readIntField(json, "type")
	if err != nil {
		return nil, err
	}

	EventTime, err := readTimeField(json, "time", time.RFC3339Nano)
	if err != nil {
		return nil, err
	}

	switch EventType(Type) {
	case EventChatShutdownType:
		reason, err := readStringField(json, "reason")
		if err != nil {
			return nil, err
		}
		return &EventChatShutdown{
			BaseEvent: BaseEvent{
				EventType: EventMessageNewType,
				EventTime: EventTime,
			},
			Reason: reason,
		}, nil
	case EventMessageNewType:
		fromID, err := readStringField(json, "from_id")
		if err != nil {
			return nil, err
		}
		body, err := readStringField(json, "body")
		if err != nil {
			return nil, err
		}
		return &EventMessageNew{
			BaseEvent: BaseEvent{
				EventType: EventMessageNewType,
				EventTime: EventTime,
			},
			FromID: fromID,
			Body:   body,
		}, nil
	case EventUserConnectedType:
		userID, err := readStringField(json, "user_id")
		if err != nil {
			return nil, err
		}
		return &EventUserConnected{
			BaseEvent: BaseEvent{
				EventType: EventUserConnectedType,
				EventTime: EventTime,
			},
			UserID: userID,
		}, nil
	case EventUserDisconnectedType:
		userID, err := readStringField(json, "user_id")
		if err != nil {
			return nil, err
		}
		return &EventUserDisconnected{
			BaseEvent: BaseEvent{
				EventType: EventUserDisconnectedType,
				EventTime: EventTime,
			},
			UserID: userID,
		}, nil
	default:
		return nil, ErrUnknownEventType
	}
}

func readIntField(json map[string]any, fieldName string) (int, error) {
	value, ok := json[fieldName]
	if !ok {
		return 0, ErrFailedToParseEvent
	}
	switch value := value.(type) {
	case int:
		return value, nil
	default:
		return 0, ErrInvalidFieldType
	}
}

func readStringField(json map[string]any, fieldName string) (string, error) {
	value, ok := json[fieldName]
	if !ok {
		return "", ErrFailedToParseEvent
	}
	switch value := value.(type) {
	case string:
		return value, nil
	default:
		return "", ErrInvalidFieldType
	}
}

func readTimeField(json map[string]any, fieldName string, layout string) (time.Time, error) {
	value, ok := json[fieldName]
	if !ok {
		return time.Time{}, ErrFailedToParseEvent
	}
	switch value.(type) {
	case string:
		break
	default:
		return time.Time{}, ErrInvalidFieldType
	}
	Time, err := time.Parse(layout, value.(string))
	if err != nil {
		return time.Time{}, ErrInvalidEventTime
	}
	return Time, nil
}
