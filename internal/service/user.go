package service

import (
	"net/http"

	"github.com/ilovepitsa/jwt_tokens/internal/entity"
	"github.com/ilovepitsa/jwt_tokens/internal/repo"
	"github.com/ilovepitsa/jwt_tokens/pkg/tokens"
)

type UserServiceInterface interface {
	SignIn(r *http.Request) (*entity.Tokens, error)
	Refresh(r *http.Request) (*entity.Tokens, error)
	CreateUser(r *http.Request) (*entity.User, error)
}

type userService struct {
	userRepo     repo.UserRepoInterface
	tokenManager tokens.TokenManager
}

func NewUserService(userRepo repo.UserRepoInterface, tokenManager tokens.TokenManager) *userService {
	return &userService{
		userRepo:     userRepo,
		tokenManager: tokenManager,
	}
}

func (us *userService) SignIn(r *http.Request) (*entity.Tokens, error) {

	return nil, nil
}

func (us *userService) Refresh(r *http.Request) (*entity.Tokens, error) {

	return nil, nil
}

func (us *userService) CreateUser(r *http.Request) (*entity.User, error) {

	return nil, nil
}
