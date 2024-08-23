package app

import (
	"errors"
	"math/rand"
	"net"
	"net/http"
	"time"

	"github.com/ilovepitsa/jwt_tokens/internal/config"
	v1 "github.com/ilovepitsa/jwt_tokens/internal/http/v1"
	"github.com/ilovepitsa/jwt_tokens/internal/repo"
	"github.com/ilovepitsa/jwt_tokens/internal/service"
	"github.com/ilovepitsa/jwt_tokens/pkg/tokens"
)

var (
	ErrReadConfig       = errors.New("cant read config")
	ErrConnectDB        = errors.New("cant reach db")
	ErrInitService      = errors.New("cant initialize service")
	ErrInitTokenManager = errors.New("cant init token manager")
)

func Run(configPath string) error {

	cfg, err := config.NewConfig(configPath)
	if err != nil {
		return errors.Join(ErrReadConfig, err)
	}

	manager, err := tokens.NewManager([]byte("must-be-secret-key"), func(b []byte) []byte {
		s := rand.NewSource(time.Now().Unix())
		r := rand.New(s)

		if _, err := r.Read(b); err != nil {
			return []byte("")
		}

		return b
	})

	if err != nil {
		return errors.Join(ErrInitTokenManager, err)
	}

	repo := repo.NewRepo(*cfg)

	dep := service.Dependencies{
		Repo:         repo,
		TokenManager: manager,
	}

	serv := service.NewServices(dep)

	handler := http.NewServeMux()

	v1.NewRouter(handler, *serv)

	err = http.ListenAndServe(net.JoinHostPort(cfg.NetworkSettings.Host, cfg.NetworkSettings.Port), handler)
	return err
}
