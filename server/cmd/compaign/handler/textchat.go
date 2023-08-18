package handler

// var (
// 	errUserIDNotFound = errors.New("user_id not passed in request")
// 	errChatIDNotFound = errors.New("chat_id not passed in request")
// )

// var wsUpgrader = websocket.Upgrader{
// 	ReadBufferSize:  1024,
// 	WriteBufferSize: 1024,
// }

// // returns { "textchat_id": "UUID" }
// func BuildTextChatCreateHandler(
// 	logger *logrus.Entry,
// 	chatMngr *service.ChatManager,
// ) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		logger := logger.WithField("r", r)
// 		logger.Info("onTextChatCreateHandler")

// 		id := uuid.NewString()
// 		err := chatMngr.CreateTextChat(id)
// 		if err != nil {
// 			logger.Error(err)
// 			w.WriteHeader(http.StatusInternalServerError)
// 			return
// 		}

// 		json.NewEncoder(w).Encode(map[string]string{
// 			"textchat_id": id,
// 		})
// 	}
// }

// func BuildTextChatDeleteHandler(
// 	logger *logrus.Entry,
// 	chatMngr *service.ChatManager,
// ) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		logger := logger.WithField("r", r)
// 		logger.Info("onTextChatDeleteHandler")

// 		chatID, ok := mux.Vars(r)["chat_id"]
// 		if !ok {
// 			logger.Error(errChatIDNotFound)
// 			w.Write([]byte(errChatIDNotFound.Error()))
// 			w.WriteHeader(http.StatusBadRequest)
// 			return
// 		}

// 		err := chatMngr.CloseTextChat(chatID)
// 		if err != nil {
// 			logger.Error(err)
// 			w.Write([]byte(err.Error()))
// 			w.WriteHeader(http.StatusBadRequest)
// 			return
// 		}

// 		w.WriteHeader(http.StatusNoContent)
// 	}
// }

// func BuildTextChatHandler(
// 	logger *logrus.Entry,
// 	chatMngr *service.ChatManager,
// ) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		logger := logger.WithField("r", r)
// 		logger.Info("onTextChatHandler")

// 		vars := mux.Vars(r)
// 		chatID, ok := vars["chat_id"]
// 		if !ok {
// 			logger.Error(errChatIDNotFound)
// 			w.Write([]byte(errChatIDNotFound.Error()))
// 			w.WriteHeader(http.StatusBadRequest)
// 			return
// 		}

// 		userID := r.FormValue("user_id")
// 		if userID == "" {
// 			logger.Error(errUserIDNotFound)
// 			w.Write([]byte(errUserIDNotFound.Error()))
// 			w.WriteHeader(http.StatusBadRequest)
// 			return
// 		}

// 		chat, err := chatMngr.GetTextChatByID(chatID)
// 		if err != nil {
// 			logger.Error(err)
// 			w.Write([]byte(err.Error()))
// 			w.WriteHeader(http.StatusNotFound)
// 			return
// 		}

// 		// TODO: check user id here

// 		conn, err := wsUpgrader.Upgrade(w, r, nil)
// 		if err != nil {
// 			logger.Error(err)
// 			return
// 		}

// 		chat.ConnectionRequested(userID, conn)
// 	}
// }
