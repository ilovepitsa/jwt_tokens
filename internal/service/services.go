package service

import (
	"github.com/ilovepitsa/jwt_tokens/internal/config"
	"github.com/ilovepitsa/jwt_tokens/internal/notification"
	"github.com/ilovepitsa/jwt_tokens/internal/repo"
	"github.com/ilovepitsa/jwt_tokens/pkg/tokens"
)

type Services struct {
	UserService UserServiceInterface
}

type Dependencies struct {
	Cfg           config.Config
	Repo          *repo.Repo
	TokenManager  tokens.TokenManager
	EmailNotifier notification.Notifier
}

func NewServices(dep Dependencies) *Services {
	services := &Services{}

	services.UserService = NewUserService(dep.Repo.UserRepo, dep.TokenManager, dep.EmailNotifier, dep.Cfg.TokensSettings.AccessTTL)

	return services
}
