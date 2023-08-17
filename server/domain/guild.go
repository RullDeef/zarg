package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/exp/slices"
)

var (
	ErrEmptyProfiles         = errors.New("empty profiles")
	ErrProfileAlreadyInGuild = errors.New("profile already in guild")
	ErrProfileNotInGuild     = errors.New("profile not in guild")

	// ErrRequestInProgress означает, что запрос на присоединение уже отправлен
	ErrRequestInProgress = errors.New("request in progress")

	ErrJoinRequestNotFound = errors.New("join request not found")
)

type GuildID string

type Guild struct {
	ID   GuildID `json:"id"`
	Name string  `json:"name"`

	// Leader - глава гильдии
	Leader *Profile `json:"leader"`

	// Participators - участники гильдии
	Participators []GuildParticipator `json:"participators"`

	// JoinRequests - запросы на присоединение
	JoinRequests []GuildJoinRequest `json:"join_requests"`

	// Activity - история походов, совершенных участниками гильдии.
	Activity []*GuildActivity `json:"activity"`
}

type GuildParticipator struct {
	Profile *Profile `json:"profile"`

	// JoinDate - время, когда данный участник присоединился к гильдии
	JoinDate time.Time `json:"join_date"`
}

type GuildJoinRequest struct {
	// Profile - от кого поступил запрос
	Profile *Profile

	// JoinRequestDate - время, когда был отправлен запрос
	JoinRequestDate time.Time

	// Message - опциональное сообщение отправляемое вместе с запросом
	Message string
}

// GuildActivity - отслеживаемая активность участников гильдии
type GuildActivity struct {
	// AcitivityType - тип активности (для отображения пользователям)
	AcitivityType string

	// Participators - участники активности
	Participators []*Profile

	// StartTime - время начала активности
	StartTime time.Time

	// Duration - длительность активности
	Duration time.Duration
}

// NewGuildFromProfiles - создает новую гильдию. Главой гильдии назначается
// первый участник в переданном массиве
func NewGuildFromProfiles(guildName string, profiles []*Profile) (*Guild, error) {
	if len(profiles) == 0 {
		return nil, ErrEmptyProfiles
	}

	participators := make([]GuildParticipator, len(profiles))
	for i, profile := range profiles {
		participators[i] = newParticipator(profile)
	}

	return &Guild{
		ID:            GuildID(uuid.New().String()),
		Name:          guildName,
		Leader:        profiles[0],
		Participators: participators,
	}, nil
}

func newParticipator(profile *Profile) GuildParticipator {
	return GuildParticipator{
		Profile:  profile,
		JoinDate: time.Now(),
	}
}

// RequestJoin - отправка запроса присоединения игрока к гильдии.
// Возвращает ошибку, если игрок уже состоит в гильдии или имеет необработанный запрос.
func (g *Guild) RequestJoin(profile *Profile, message string) error {
	if g.hasProfileInParticipators(profile) {
		return ErrProfileAlreadyInGuild
	}
	if g.hasProfileInJoinRequests(profile) {
		return ErrRequestInProgress
	}

	g.JoinRequests = append(g.JoinRequests, newJoinRequest(profile, message))
	return nil
}

func (g *Guild) hasProfileInParticipators(profile *Profile) bool {
	for _, participator := range g.Participators {
		if participator.Profile == profile {
			return true
		}
	}
	return false
}

func (g *Guild) hasProfileInJoinRequests(profile *Profile) bool {
	for _, request := range g.JoinRequests {
		if request.Profile == profile {
			return true
		}
	}
	return false
}

func newJoinRequest(profile *Profile, message string) GuildJoinRequest {
	return GuildJoinRequest{
		Profile:         profile,
		JoinRequestDate: time.Now(),
		Message:         message,
	}
}

// AcceptJoinRequest - принять запрос на присоединение игрока к гильдии
func (g *Guild) AcceptJoinRequest(profile *Profile) error {
	if g.hasProfileInParticipators(profile) {
		return ErrProfileAlreadyInGuild
	}

	if !g.removeJoinRequest(profile) {
		return ErrJoinRequestNotFound
	}

	g.Participators = append(g.Participators, newParticipator(profile))
	return nil
}

// RejectJoinRequest - отклонить запрос на присоединение игрока к гильдии.
// При этом игрок может повторно отправить запрос на присоединение
func (g *Guild) RejectJoinRequest(profile *Profile) error {
	if g.hasProfileInParticipators(profile) {
		return ErrProfileAlreadyInGuild
	}

	if !g.removeJoinRequest(profile) {
		return ErrJoinRequestNotFound
	}

	return nil
}

func (g *Guild) removeJoinRequest(profile *Profile) bool {
	for i, request := range g.JoinRequests {
		if request.Profile == profile {
			g.JoinRequests = append(g.JoinRequests[:i], g.JoinRequests[i+1:]...)
			return true
		}
	}
	return false
}

// GetActivityForProfile - получить историю активности для участника гильдии
func (g *Guild) GetActivityForProfile(profile *Profile) ([]*GuildActivity, error) {
	if !g.hasProfileInParticipators(profile) {
		return nil, ErrProfileNotInGuild
	}

	var activities []*GuildActivity
	for _, activity := range g.Activity {
		if slices.Contains(activity.Participators, profile) {
			activities = append(activities, activity)
		}
		activities = append(activities, activity)
	}
	return activities, nil
}
