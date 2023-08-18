package textws

import (
	"container/list"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

const defaultCloseReason = "chat was closed by server"

type TextChat_OLD struct {
	chatID string
	logger *logrus.Entry

	clients *list.List
	mutex   sync.Mutex
}

type client struct {
	userID string
	conn   *websocket.Conn
}

func NewTextChat(logger *logrus.Entry) *TextChat_OLD {
	return &TextChat_OLD{
		chatID:  uuid.NewString(),
		logger:  logger,
		clients: list.New(),
		mutex:   sync.Mutex{},
	}
}

func (chat *TextChat_OLD) Close() {
	event := newEventChatShutdown(defaultCloseReason)
	chat.broadcastEvent(event)

	chat.mutex.Lock()
	defer chat.mutex.Unlock()

	for node := chat.clients.Front(); node != nil; node = node.Next() {
		client := node.Value.(*client)
		err := client.conn.Close()
		if err != nil {
			chat.logger.WithField("userID", client.userID).Error(err)
		}
	}

	chat.clients.Init()
}

func (chat *TextChat_OLD) ConnectionRequested(userID string, conn *websocket.Conn) {
	chat.logger.WithField("userID", userID).Info("ConnectionRequested")

	client := &client{
		userID: userID,
		conn:   conn,
	}
	chat.connectClient(client)

	defaultCloseHandler := conn.CloseHandler()
	conn.SetCloseHandler(func(code int, text string) error {
		chat.disconnectClient(client)
		return defaultCloseHandler(code, text)
	})

	go chat.receiverLoop(client)
}

func (chat *TextChat_OLD) receiverLoop(client *client) {
	chat.logger.WithField("userID", client.userID).Info("receiverLoop")

	for {
		msgType, bytes, err := client.conn.ReadMessage()
		if err != nil {
			break
		} else if msgType != websocket.TextMessage {
			chat.logger.WithField("userID", client.userID).
				WithField("msgType", msgType).
				Warn("received message with invalid type")
			continue
		}

		// parse incomming event
		msg := string(bytes)
		event := newEventMessageNew(client.userID, msg)
		chat.broadcastEvent(event)
	}
}

func (chat *TextChat_OLD) connectClient(clien *client) {
	chat.logger.WithField("userID", clien.userID).Info("connectClient")

	chat.mutex.Lock()
	defer chat.mutex.Unlock()

	// send everyone event with new connected user
	event := newEventUserConnected(clien.userID)
	chat.broadcastEvent(event)
	chat.clients.PushBack(clien)
}

func (chat *TextChat_OLD) disconnectClient(client *client) {
	chat.logger.WithField("userID", client.userID).Info("disconnectClient")

	chat.mutex.Lock()
	defer chat.mutex.Unlock()

	node := chat.clients.Front()
	for node != nil && node.Value != client {
		node = node.Next()
	}

	if node == nil {
		chat.logger.WithField("userID", client.userID).Error("failed to disconnect non-existent client")
		return
	}

	event := NewEventUserDisconnected(client.userID)
	chat.broadcastEvent(event)
	chat.clients.Remove(node)
}

func (chat *TextChat_OLD) broadcastEvent(event Event) {
	chat.logger.WithField("event", event).Info("broadcastEvent")

	for node := chat.clients.Front(); node != nil; node = node.Next() {
		conn := node.Value.(*client).conn
		go chat.sendEvent(conn, event)
	}
}

func (chat *TextChat_OLD) sendEvent(conn *websocket.Conn, event Event) {
	if err := conn.WriteJSON(event); err != nil {
		chat.logger.WithField("event", event).Error(err)
	}
}
