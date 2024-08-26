package tokens

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
)

type RefreshGenerator func([]byte) []byte

type TokenManager interface {
	NewJWT(userId uint32, user_ip string, ttl time.Duration) (string, error)
	Parse(inputToken string) (*CustomerInfo, error)
}

type Manager struct {
	secret    []byte
	generator RefreshGenerator
}

type CustomerInfo struct {
	*jwt.StandardClaims
	Ip string
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

func (m *Manager) NewJWT(userId uint32, user_ip string, ttl time.Duration) (string, error) {
	jwt.GetSigningMethod("")
	token := jwt.New(jwt.SigningMethodHS512)
	token.Claims = &CustomerInfo{
		&jwt.StandardClaims{
			ExpiresAt: time.Now().Add(ttl).Unix(),
			Subject:   strconv.FormatUint(uint64(userId), 10),
		},
		user_ip,
	}
	// token := jwt.NewWithClaims(jwt.SigningMethodHS512, })

	return token.SignedString(m.secret)
}

func (m *Manager) Parse(inputToken string) (*CustomerInfo, error) {

	token, err := jwt.ParseWithClaims(inputToken, &CustomerInfo{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return m.secret, nil

	})

	if err != nil {
		return &CustomerInfo{}, err
	}

	claims, ok := token.Claims.(*CustomerInfo)
	if !ok {
		return &CustomerInfo{}, fmt.Errorf("error get user claims from token")
	}

	return claims, nil
}
