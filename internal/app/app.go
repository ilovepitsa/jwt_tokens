package app

import (
	"errors"
	"log"
	"math/rand"
	"net"
	"net/http"
	"time"

	"github.com/ilovepitsa/jwt_tokens/internal/config"
	v1 "github.com/ilovepitsa/jwt_tokens/internal/http/v1"
	"github.com/ilovepitsa/jwt_tokens/internal/notification"
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

	log.Println("Reading config...")
	cfg, err := config.NewConfig(configPath)
	if err != nil {
		return errors.Join(ErrReadConfig, err)
	}

	log.Println("Initialize token manager...")
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

	log.Println("Initialize repo...")
	repo, err := repo.NewRepo(*cfg)
	if err != nil {
		log.Println("error connect db")
		return err
	}

	dep := service.Dependencies{
		Cfg:           *cfg,
		Repo:          repo,
		TokenManager:  manager,
		EmailNotifier: notification.NewEmailNotificator(),
	}

	log.Println("Initialize services...")
	serv := service.NewServices(dep)

	handler := http.NewServeMux()

	log.Println("Initialize router...")
	v1.NewRouter(handler, *serv)

	host := net.JoinHostPort(cfg.NetworkSettings.Host, cfg.NetworkSettings.Port)
	log.Println("Starting server on ", host)
	err = http.ListenAndServe(host, handler)
	return err
}
