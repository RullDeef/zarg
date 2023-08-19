package lobby

import (
	"context"
	"server/domain"
	"server/internal/modules/auth"
)

type Controller struct {
	authManager *auth.AuthManager
	lobby       *Lobby
}

func NewController(
	authManager *auth.AuthManager,
	lobby *Lobby,
) *Controller {
	return &Controller{
		authManager: authManager,
		lobby:       lobby,
	}
}

// AcceptJoinRequest - обработать запрос на участие в походе от пользователя
// с переданным токеном. Если токен пуст - пользователь считается анонимным.
// Возвращает идентификатор пользователя (нового, если аноним)
func (c *Controller) AcceptJoinRequest(mode string, authToken string) (domain.ProfileID, error) {
	if profile, err := c.auth(authToken); err != nil {
		return "", err
	} else if err = c.lobby.AddRequest(ParticipationMode(mode), profile); err != nil {
		return "", err
	} else {
		return profile.ID, err
	}
}

func (c *Controller) CancelRequest(profileID domain.ProfileID) {
	c.lobby.CancelRequestByID(profileID)
}

// WaitJoin - ожидает формирования команды (блокирующий вызов).
// Если пользователь не находится ни в одной очереди - сразу возвращает ошибку.
func (c *Controller) WaitJoin(ctx context.Context, profileID domain.ProfileID) (domain.CompaignID, error) {
	return c.lobby.WaitJoin(ctx, profileID)
}

func (c *Controller) auth(authToken string) (*domain.Profile, error) {
	if authToken != "" {
		return c.authManager.ValidateToken(authToken)
	} else {
		return c.authManager.AuthAnonymous()
	}
}
