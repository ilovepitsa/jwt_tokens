package repo

import (
	"github.com/ilovepitsa/jwt_tokens/internal/config"
)

type Repo struct {
	UserRepo UserRepoInterface
}

func NewRepo(cfg config.Config) *Repo {

	return &Repo{}
}
