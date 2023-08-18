package handler

// func BuildProfileGet(logger *logrus.Entry, repo *repository.Profile) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		logger := logger.WithField("r", r)
// 		logger.Info("profileGet")

// 		vars := mux.Vars(r)
// 		profileID := vars["profile_id"]

// 		profile, err := repo.GetByID(profileID)
// 		if err != nil {
// 			logger.Error("not found")
// 			w.WriteHeader(http.StatusNotFound)
// 		} else {
// 			json.NewEncoder(w).Encode(profile)
// 		}
// 	}
// }

// func BuildProfileCreateAnonymous(logger *logrus.Entry, repo *repository.Profile) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		logger := logger.WithField("r", r)
// 		logger.Info("profileCreateAnonymous")

// 		profile := smodel.NewAnonymousProfile()
// 		err := repo.Insert(profile)
// 		if err != nil {
// 			logger.WithField("profile", profile).Error(err)
// 			w.WriteHeader(http.StatusInternalServerError)
// 		} else {
// 			json.NewEncoder(w).Encode(profile)
// 		}
// 	}
// }
