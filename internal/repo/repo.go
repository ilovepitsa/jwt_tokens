package repo

import (
	"context"

	"github.com/ilovepitsa/jwt_tokens/internal/config"
)

type Repo struct {
	UserRepo UserRepoInterface
}

func NewRepo(cfg config.Config) (*Repo, error) {

	repo := &Repo{}
	uRepo, err := NewUserRepo(context.TODO(), cfg)
	if err != nil {
		return nil, err
	}
	repo.UserRepo = uRepo

	return repo, nil
}
