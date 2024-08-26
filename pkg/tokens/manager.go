package tokens

import (
	"bytes"
	b64 "encoding/base64"
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
	ParseRefreshToken(inputToken string) (*RefreshInfo, error)
	NewRefreshToken(userId uint32, user_ip string, ttl time.Duration) (string, string, error)
}

type Manager struct {
	secret    []byte
	generator RefreshGenerator
}

type CustomerInfo struct {
	*jwt.StandardClaims
	Ip string
}

type RefreshInfo struct {
	*CustomerInfo
	RefreshToken string
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

func (m *Manager) NewRefreshToken(userId uint32, user_ip string, ttl time.Duration) (string, string, error) {
	token := jwt.New(jwt.SigningMethodHS512)
	rawToken, err := m.newRawRefreshToken()
	if err != nil {
		return "", "", err
	}
	token.Claims = &RefreshInfo{
		&CustomerInfo{
			&jwt.StandardClaims{
				ExpiresAt: time.Now().Add(ttl).Unix(),
				Subject:   strconv.FormatUint(uint64(userId), 10),
			}, user_ip,
		}, rawToken,
	}
	// token := jwt.NewWithClaims(jwt.SigningMethodHS512, })

	refreshJWT, err := token.SignedString(m.secret)

	return refreshJWT, rawToken, err

}

func (m *Manager) newRawRefreshToken() (string, error) {
	b := make([]byte, 32)

	b = m.generator(b)
	enc := b64.URLEncoding.EncodeToString(b)
	return enc, nil
}

func (m *Manager) ParseRefreshToken(inputToken string) (*RefreshInfo, error) {
	token, err := jwt.ParseWithClaims(inputToken, &RefreshInfo{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return m.secret, nil
	})

	if err != nil {
		return &RefreshInfo{}, err
	}

	claims, ok := token.Claims.(*RefreshInfo)
	if !ok {
		return &RefreshInfo{}, fmt.Errorf("error get user claims from token")
	}

	return claims, nil
}
