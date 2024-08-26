package service

import (
	"context"
	"time"

	"github.com/ilovepitsa/jwt_tokens/internal/entity"
	"github.com/ilovepitsa/jwt_tokens/internal/repo"
	"github.com/ilovepitsa/jwt_tokens/pkg/tokens"
)

type UserServiceInterface interface {
	SignIn(user_id uint32) (*entity.Tokens, error)
	Refresh(refresh_toker string) (*entity.Tokens, error)
	CreateUser() (*entity.User, error)
}

type userService struct {
	userRepo     repo.UserRepoInterface
	tokenManager tokens.TokenManager
	tokenTTL     time.Duration
}

func NewUserService(userRepo repo.UserRepoInterface, tokenManager tokens.TokenManager, tokenTTL time.Duration) *userService {
	return &userService{
		userRepo:     userRepo,
		tokenManager: tokenManager,
		tokenTTL:     tokenTTL,
	}
}

func (us *userService) createSession(userId uint32) (*entity.Tokens, error) {
	access, err := us.tokenManager.NewJWT(userId, us.tokenTTL)
	if err != nil {
		return nil, err
	}
	refresh, err := us.tokenManager.NewRefreshToken()
	if err != nil {
		return nil, err
	}

	err = us.userRepo.CreateSession(context.TODO(), userId, refresh)
	if err != nil {
		return nil, err
	}

	return &entity.Tokens{
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}

func (us *userService) SignIn(user_id uint32) (*entity.Tokens, error) {

	exist, err := us.userRepo.CheckUserExist(context.TODO(), user_id)
	if err == repo.ErrUserNotExist && !exist {
		return nil, err
	}
	return us.createSession(user_id)
}

func (us *userService) Refresh(refresh_token string) (*entity.Tokens, error) {

	user, err := us.userRepo.GetByRefreshToken(context.TODO(), refresh_token)
	if err != nil {
		return nil, err
	}

	return us.createSession(user)
}

func (us *userService) CreateUser() (*entity.User, error) {
	id, err := us.userRepo.CreateUser(context.TODO())
	if err != nil {
		return nil, err
	}
	return &entity.User{Id: id}, nil
}
