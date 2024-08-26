package service

import (
	"context"
	"fmt"
	"time"

	"github.com/ilovepitsa/jwt_tokens/internal/entity"
	"github.com/ilovepitsa/jwt_tokens/internal/notification"
	"github.com/ilovepitsa/jwt_tokens/internal/repo"
	"github.com/ilovepitsa/jwt_tokens/pkg/tokens"
)

type UserServiceInterface interface {
	SignIn(user_id uint32, user_ip string) (*entity.Tokens, error)
	Refresh(refresh_toker string, user_ip string) (*entity.Tokens, error)
	CreateUser() (*entity.User, error)
}

type userService struct {
	userRepo      repo.UserRepoInterface
	tokenManager  tokens.TokenManager
	tokenTTL      time.Duration
	emailNotifier notification.Notifier
}

func NewUserService(userRepo repo.UserRepoInterface, tokenManager tokens.TokenManager, emailNotifier notification.Notifier, tokenTTL time.Duration) *userService {
	return &userService{
		userRepo:      userRepo,
		tokenManager:  tokenManager,
		tokenTTL:      tokenTTL,
		emailNotifier: emailNotifier,
	}
}

func (us *userService) createSession(userId uint32, user_ip string) (*entity.Tokens, error) {
	access, err := us.tokenManager.NewJWT(userId, user_ip, us.tokenTTL)
	if err != nil {
		return nil, err
	}
	refresh, err := us.tokenManager.NewRefreshToken()
	if err != nil {
		return nil, err
	}

	err = us.userRepo.CreateSession(context.TODO(), userId, refresh, user_ip)
	if err != nil {
		return nil, err
	}

	return &entity.Tokens{
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}

func (us *userService) SignIn(user_id uint32, user_ip string) (*entity.Tokens, error) {

	exist, err := us.userRepo.CheckUserExist(context.TODO(), user_id)
	if err == repo.ErrUserNotExist && !exist {
		return nil, err
	}
	return us.createSession(user_id, user_ip)
}

func (us *userService) Refresh(refresh_token string, user_ip string) (*entity.Tokens, error) {

	user, prev_ip, err := us.userRepo.GetByRefreshToken(context.TODO(), refresh_token)
	if err != nil {
		return nil, err
	}
	if prev_ip != user_ip {
		us.emailNotifier.Notificate(fmt.Sprintf("entered from new device: %s", user_ip))
	}

	return us.createSession(user, user_ip)
}

func (us *userService) CreateUser() (*entity.User, error) {
	id, err := us.userRepo.CreateUser(context.TODO())
	if err != nil {
		return nil, err
	}
	return &entity.User{Id: id}, nil
}
