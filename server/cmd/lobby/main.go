package main

import (
	"fmt"
	"server/internal/modules/auth"
	"server/internal/modules/lobby"
	"server/internal/modules/logger"
	"server/internal/modules/server"

	"go.uber.org/fx"
)

const defaultPort = 4667

func main() {
	app := fx.New(
		logger.Module,
		lobby.Module,
		auth.Module,
		server.Module,
		fx.Invoke(func(s *server.Server) {
			s.Run(fmt.Sprintf(":%d", defaultPort))
		}),
	)

	app.Run()
}

// func main() {

// 	log := logrus.New()
// 	// log.SetReportCaller(true)
// 	logger := log.WithField("microservice", "lobby")

// 	profileRepo := repository.NewProfile(logger)
// 	participator := service.NewParticipator(logger, buildCompaignFactory(logger))

// 	router := mux.NewRouter()
// 	profileRouter := router.PathPrefix("/profiles").Subrouter()
// 	profileRouter.HandleFunc("/new", handler.BuildProfileCreateAnonymous(logger, profileRepo)).Methods("POST")
// 	profileRouter.HandleFunc("/{profile_id}", handler.BuildProfileGet(logger, profileRepo)).Methods("GET")

// 	router.HandleFunc("/compaigns/new", handler.BuildOnNewCompaignRequest(logger, participator))

// 	http.Handle("/", router)

// 	if err := http.ListenAndServe(fmt.Sprintf(":%d", defaultPort), nil); err != nil {
// 		logger.Error(err)
// 		panic(err)
// 	}
// }

// func buildCompaignFactory(logger *logrus.Entry) service.CompaignFactoryFunc {
// 	return func(team []service.ParticipationRequest) (map[string]any, error) {
// 		playerIDs := make([]string, 0, len(team))
// 		for _, req := range team {
// 			playerIDs = append(playerIDs, req.ProfileID)
// 		}

// 		textChatID, err := createTextChat(logger)
// 		if err != nil {
// 			logger.Error(err)
// 			return nil, err
// 		}

// 		compaignID, err := createCompaign()
// 		if err != nil {
// 			logger.Error(err)
// 			deleteTextChat(textChatID)
// 			return nil, err
// 		}

// 		message := map[string]any{
// 			"compaign_id": compaignID,
// 			"textchat_id": textChatID,
// 			"players":     playerIDs,
// 		}

// 		return message, nil
// 	}
// }

// func createTextChat(logger *logrus.Entry) (string, error) {
// 	resp, err := http.Post("http://localhost:4668/textchats", "none", nil)
// 	if err != nil {
// 		return "", err
// 	}

// 	// buffer, err := io.ReadAll(resp.Body)
// 	// logger.Warn(string(buffer), err)

// 	data := make(map[string]any)
// 	err = json.NewDecoder(resp.Body).Decode(&data)
// 	if err != nil {
// 		return "", err
// 	}

// 	textChatID, ok := data["textchat_id"]
// 	if !ok {
// 		return "", errors.New("invalid response")
// 	}
// 	return textChatID.(string), nil
// }

// func deleteTextChat(id string) {
// 	http.NewRequest(http.MethodDelete, "http://localhost:4668/textchats/"+id, nil)
// }

// func createCompaign() (string, error) {
// 	resp, err := http.Post("http://localhost:4668/compaign", "none", nil)
// 	if err != nil {
// 		return "", err
// 	}

// 	data := make(map[string]string)
// 	err = json.NewDecoder(resp.Body).Decode(&data)
// 	if err != nil {
// 		return "", err
// 	}

// 	compaignID, ok := data["compaign_id"] // TODO: fix source
// 	if !ok {
// 		return "", errors.New("invalid response")
// 	}
// 	return compaignID, nil
// }
