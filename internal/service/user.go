package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/ilovepitsa/jwt_tokens/internal/entity"
	"github.com/ilovepitsa/jwt_tokens/internal/notification"
	"github.com/ilovepitsa/jwt_tokens/internal/repo"
	"github.com/ilovepitsa/jwt_tokens/pkg/tokens"
	"golang.org/x/crypto/bcrypt"
)

type UserServiceInterface interface {
	SignIn(user_id uint32, user_ip string) (*entity.Tokens, error)
	Refresh(refresh_toker string, user_ip string) (*entity.Tokens, error)
	CreateUser() (*entity.User, error)
}

type userService struct {
	userRepo        repo.UserRepoInterface
	tokenManager    tokens.TokenManager
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
	emailNotifier   notification.Notifier
}

var (
	ErrRefreshNotEqual = errors.New("refresh token not equal")
)

func NewUserService(userRepo repo.UserRepoInterface, tokenManager tokens.TokenManager, emailNotifier notification.Notifier, accessTokenTTL time.Duration, refreshTokenTTL time.Duration) *userService {
	return &userService{
		userRepo:        userRepo,
		tokenManager:    tokenManager,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
		emailNotifier:   emailNotifier,
	}
}

// В РИДМИ ПИШУ ПРО СВОИ ЛОГИЧЕСКИЕ НЕСОСТЫКОВКИ, ЧТО БКРИПТ ВСЕГДА РАЗНЫЙ
// ПОЭТОМУ БЫЛА ВЫДУМАНА ТАКАЯ СХЕМА ИЗ КОСТЫЛЕЙ
func (us *userService) createSession(userId uint32, user_ip string) (*entity.Tokens, error) {
	access, err := us.tokenManager.NewJWT(userId, user_ip, us.accessTokenTTL)
	if err != nil {
		return nil, err
	}
	refreshJWT, refreshRaw, err := us.tokenManager.NewRefreshToken(userId, user_ip, us.refreshTokenTTL)
	if err != nil {
		return nil, err
	}

	log.Println("refresh token raw: ", refreshRaw)
	err = us.userRepo.CreateSession(context.TODO(), userId, refreshRaw, user_ip)
	if err != nil {
		return nil, err
	}

	return &entity.Tokens{
		AccessToken:  access,
		RefreshToken: refreshJWT,
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

	info, err := us.tokenManager.ParseRefreshToken(refresh_token)
	if err != nil {
		return nil, err
	}
	userId, err := strconv.ParseUint(info.Subject, 10, 32)
	if err != nil {
		return nil, err
	}

	old_token, prev_ip, err := us.userRepo.GetRefreshByInfo(context.TODO(), uint32(userId))

	if err != nil {
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword(old_token, []byte(info.RefreshToken)); err != nil {
		return nil, ErrRefreshNotEqual
	}

	if prev_ip != user_ip {
		us.emailNotifier.Notificate(fmt.Sprintf("entered from new device: %s", user_ip))
	}

	return us.createSession(uint32(userId), user_ip)
}

func (us *userService) CreateUser() (*entity.User, error) {
	id, err := us.userRepo.CreateUser(context.TODO())
	if err != nil {
		return nil, err
	}
	return &entity.User{Id: id}, nil
}
