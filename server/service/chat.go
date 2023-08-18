package service

import (
	"errors"
	"sync"

	"github.com/sirupsen/logrus"
)

var (
	ErrChatExists   = errors.New("chat with given id already exists")
	ErrChatNotFound = errors.New("chat with given id not found")
)

type ChatManager struct {
	logger *logrus.Entry
	// textChats map[string]*textchat.TextChat
	mutex sync.RWMutex
}

func NewChatManager(logger *logrus.Entry) *ChatManager {
	return &ChatManager{
		logger: logger,
		// textChats: make(map[string]*textchat.TextChat),
		mutex: sync.RWMutex{},
	}
}

func (mngr *ChatManager) CreateTextChat(id string) error {
	mngr.logger.WithField("id", id).Info("CreateTextChat")

	mngr.mutex.Lock()
	defer mngr.mutex.Unlock()

	// _, ok := mngr.textChats[id]
	// if ok {
	// 	mngr.logger.WithField("id", id).Error(ErrChatExists)
	// 	return ErrChatExists
	// } else {
	// 	mngr.textChats[id] = textchat.NewTextChat(mngr.logger)
	// 	mngr.logger.WithField("id", id).Info("textchat created successfully")
	// 	return nil
	// }
	return nil
}

func (mngr *ChatManager) CloseTextChat(id string) error {
	mngr.logger.WithField("id", id).Info("CloseTextChat")

	mngr.mutex.Lock()
	defer mngr.mutex.Unlock()

	// chat, ok := mngr.textChats[id]
	// if !ok {
	// 	mngr.logger.WithField("id", id).Error(ErrChatNotFound)
	// 	return ErrChatNotFound
	// }

	// chat.Close()
	// delete(mngr.textChats, id)
	return nil
}

// func (mngr *ChatManager) GetTextChatByID(id string) (*textchat.TextChat, error) {
// 	mngr.logger.WithField("id", id).Info("GetTextChatByID")

// 	mngr.mutex.RLock()
// 	defer mngr.mutex.RUnlock()

// 	chat, ok := mngr.textChats[id]
// 	if !ok {
// 		mngr.logger.WithField("id", id).Error(ErrChatNotFound)
// 		return nil, ErrChatNotFound
// 	}

// 	return chat, nil
// }
