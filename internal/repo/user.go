package repo

import (
	"context"
	"fmt"

	"github.com/ilovepitsa/jwt_tokens/internal/config"
	"github.com/jackc/pgx/v5"
)

type UserRepoInterface interface {
	CreateUser() (uint32, error)
	GetByRefreshToken(string) (uint32, error)
}

type userRepo struct {
	conn *pgx.Conn
}

func NewUserRepo(ctx context.Context, cfg config.Config) (*userRepo, error) {
	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.PostgresSettings.User, cfg.PostgresSettings.Password,
		cfg.PostgresSettings.Host, cfg.PostgresSettings.Port, cfg.PostgresSettings.DB)

	conn, err := pgx.Connect(ctx, url)
	if err != nil {
		return nil, err
	}
	err = conn.Ping(ctx)
	if err != nil {
		return nil, err
	}

	return &userRepo{conn: conn}, nil

}

func (r *userRepo) CreateUser() (uint32, error) {
	return 0, nil
}

func (r *userRepo) GetByRefreshToken(refreshToken string) (uint32, error) {

	return 0, nil
}
