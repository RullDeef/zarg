package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Module("auth",
	fx.Provide(func(logger *zap.SugaredLogger) *AuthManager {
		return New(logger, "secret", jwt.SigningMethodEdDSA)
	}),
)
