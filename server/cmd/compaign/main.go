package main

// const defaultPort = 4668

// var logger *logrus.Entry

// func main() {
// 	log := logrus.New()
// 	// log.SetReportCaller(true)
// 	logger = log.WithField("microservice", "compaign")

// 	chatMngr := service.NewChatManager(logger)

// 	router := mux.NewRouter()

// 	compaignRouter := router.PathPrefix("/compaign").Subrouter()

// 	compaignRouter.
// 		HandleFunc("", handler.BuildCreateCompaignHandler(logger)).
// 		Methods("POST")
// 	compaignRouter.
// 		HandleFunc("/{id}", handler.GetCompaign).
// 		Methods("GET")

// 	textChatsRouter := router.PathPrefix("/textchats").Subrouter()
// 	// creates brand new text chat
// 	textChatsRouter.
// 		HandleFunc("", handler.BuildTextChatCreateHandler(logger, chatMngr)).
// 		Methods("POST")
// 	// connects user to specific text chat
// 	textChatsRouter.Path("/{chat_id}").
// 		Queries("user_id", "{user_id:.*}").
// 		HandlerFunc(handler.BuildTextChatHandler(logger, chatMngr))
// 	// deletes text chat (DANGEROUS!!) // TODO: add confirmation token of some kind to allow only server execute request
// 	textChatsRouter.Path("/{chat_id}").
// 		HandlerFunc(handler.BuildTextChatDeleteHandler(logger, chatMngr)).
// 		Methods("DELETE")

// 	http.Handle("/", router)
// 	err := http.ListenAndServe(fmt.Sprintf(":%d", defaultPort), nil)
// 	if err != nil {
// 		logger.Error(err)
// 		panic(err)
// 	}
// }
