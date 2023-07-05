package repository

import (
	"container/list"
	"errors"

	"server/model"

	"github.com/sirupsen/logrus"
)

var (
	errProfileNotFound = errors.New("profile not found")
	errProfileExists   = errors.New("profile already exists")
)

type Profile struct {
	logger   *logrus.Entry
	profiles *list.List
}

func NewProfile(logger *logrus.Entry) *Profile {
	return &Profile{
		logger:   logger,
		profiles: list.New(),
	}
}

func (repo *Profile) Insert(profile *model.Profile) error {
	repo.logger.WithField("profile", profile).Info("Insert")

	if _, err := repo.GetByID(profile.ID); err != errProfileNotFound {
		repo.logger.WithField("profile", profile).Info(errProfileExists)
		return errProfileExists
	}

	repo.profiles.PushBack(profile)
	return nil
}

func (repo *Profile) GetByID(id string) (*model.Profile, error) {
	repo.logger.WithField("id", id).Info("GetByID")

	for node := repo.profiles.Front(); node != nil; node = node.Next() {
		profile := node.Value.(*model.Profile)
		if profile.ID == id {
			return profile, nil
		}
	}

	repo.logger.WithField("id", id).Info(errProfileNotFound)
	return nil, errProfileNotFound
}
