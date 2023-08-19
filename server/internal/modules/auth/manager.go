package auth

import (
	"errors"
	"fmt"
	"server/domain"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

var (
	ErrInvalidClaims   = errors.New("invalid claims")
	ErrProfileNotFound = errors.New("profile not found")
)

type AuthManager struct {
	logger           *zap.SugaredLogger
	jwtSecret        string
	jwtSigningMethod jwt.SigningMethod
	profiles         map[domain.ProfileID]*domain.Profile
}

func New(
	logger *zap.SugaredLogger,
	jwtSecret string,
	jwtSigningMethod jwt.SigningMethod,
) *AuthManager {
	return &AuthManager{
		logger:           logger,
		jwtSecret:        jwtSecret,
		jwtSigningMethod: jwtSigningMethod,
		profiles:         make(map[domain.ProfileID]*domain.Profile),
	}
}

// GetToken - создает токен для переданного профиля.
func (am *AuthManager) GetToken(profile *domain.Profile) (string, error) {
	am.logger.Infow("GetToken", "profile", profile)

	jwtToken := jwt.NewWithClaims(am.jwtSigningMethod, jwt.MapClaims{
		"id":  profile.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenStr, err := jwtToken.SignedString(am.jwtSecret)
	if err != nil {
		return "", err
	}

	am.profiles[profile.ID] = profile
	return tokenStr, nil
}

// ValidateToken - проверяет переданный токен и
// возвращает профиль, который ему соответствует.
func (am *AuthManager) ValidateToken(token string) (*domain.Profile, error) {
	am.logger.Infow("ValidateToken", "token", token)

	if jwtToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if token.Method != am.jwtSigningMethod {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return am.jwtSecret, nil
	}); err != nil {
		return nil, err
	} else if claims, ok := jwtToken.Claims.(jwt.MapClaims); !ok {
		return nil, ErrInvalidClaims
	} else if profile, ok := am.profiles[domain.ProfileID(claims["id"].(string))]; !ok {
		return nil, ErrProfileNotFound
	} else {
		return profile, nil
	}
}

// AuthAnonymous - создает анонимного пользователя
func (am *AuthManager) AuthAnonymous() (*domain.Profile, error) {
	return domain.NewAnonymousProfile(), nil
}
