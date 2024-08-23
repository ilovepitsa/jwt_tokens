package tokens

import (
	"bytes"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

type RefreshGenerator func([]byte) []byte

type TokenManager interface {
	NewJWT(userId string, ttl time.Duration) (string, error)
	Parse(accessToken string) (string, error)
	NewRefreshToken() (string, error)
}

type Manager struct {
	secret    []byte
	generator RefreshGenerator
}

func NewManager(secretKey []byte, generator RefreshGenerator) (*Manager, error) {
	if bytes.Equal([]byte{}, secretKey) {
		return nil, errors.New("empty secret key")
	}

	return &Manager{
		secret:    secretKey,
		generator: generator,
	}, nil
}

func (m *Manager) NewJWT(userId string, ttl time.Duration) (string, error) {
	jwt.GetSigningMethod("")
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(ttl).Unix(),
		Subject:   userId,
	})

	return token.SignedString(m.secret)
}

func (m *Manager) Parse(accessToken string) (string, error) {
	token, err := jwt.Parse(accessToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return m.secret, nil
	})

	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("error get user claims from token")
	}

	return claims["sub"].(string), nil
}

func (m *Manager) NewRefreshToken() (string, error) {
	b := make([]byte, 32)

	b = m.generator(b)

	return fmt.Sprintf("%x", b), nil
}
