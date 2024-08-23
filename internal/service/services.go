package service

import (
	"github.com/ilovepitsa/jwt_tokens/internal/repo"
	"github.com/ilovepitsa/jwt_tokens/pkg/tokens"
)

type Services struct {
	UserService UserServiceInterface
}

type Dependencies struct {
	Repo         *repo.Repo
	TokenManager tokens.TokenManager
}

func NewServices(dep Dependencies) *Services {
	services := &Services{}
	services.UserService = NewUserService(dep.Repo.UserRepo, dep.TokenManager)

	return services
}
