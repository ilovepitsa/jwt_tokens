package repo

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/ilovepitsa/jwt_tokens/internal/config"
	"github.com/jackc/pgx/v5"
)

var (
	ErrUserNotExist = errors.New("user not exist")
)

type UserRepoInterface interface {
	CreateUser(context.Context) (uint32, error)
	GetByRefreshToken(context.Context, string) (uint32, string, error)
	CheckUserExist(context.Context, uint32) (bool, error)
	CreateSession(ctx context.Context, userId uint32, refreshToken string, userIp string) error
}

type userSql struct {
	Id sql.NullInt32
}

type userRepo struct {
	conn       *pgx.Conn
	refreshTTL time.Duration
}

func NewUserRepo(ctx context.Context, cfg config.Config) (*userRepo, error) {
	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.PostgresSettings.User, cfg.PostgresSettings.Password,
		cfg.PostgresSettings.Host, cfg.PostgresSettings.Port, cfg.PostgresSettings.DB)
	var conn *pgx.Conn
	var err error

	for attemptsLeft := cfg.PostgresSettings.ConnectionAttempts; attemptsLeft > 0; attemptsLeft-- {
		log.Println("Trying to connect to database...")
		log.Println("Attempts left: ", attemptsLeft)
		conn, err = pgx.Connect(ctx, url)
		if err == nil {
			break
		}
		time.Sleep(cfg.PostgresSettings.ConntectTimeout)
	}

	if err != nil {
		return nil, err
	}
	err = conn.Ping(ctx)
	if err != nil {
		return nil, err
	}

	return &userRepo{
		conn:       conn,
		refreshTTL: cfg.TokensSettings.RefreshTTL}, nil

}

func (r *userRepo) CreateUser(ctx context.Context) (uint32, error) {
	trans, err := r.conn.Begin(ctx)
	if err != nil {
		trans.Rollback(ctx)
		return 0, err
	}
	var id uint32
	err = trans.QueryRow(ctx, "insert into users DEFAULT values RETURNING id;").Scan(&id)
	if err != nil {
		trans.Rollback(ctx)
		return 0, err
	}
	trans.Commit(ctx)
	return id, nil
}

func (r *userRepo) GetByRefreshToken(ctx context.Context, refreshToken string) (uint32, string, error) {

	trans, err := r.conn.Begin(ctx)
	if err != nil {
		trans.Rollback(ctx)
		return 0, "", err
	}
	var user_id int
	var prev_ip string

	h := sha256.New()
	h.Write([]byte(refreshToken))
	hashed_token := h.Sum(nil)

	err = trans.QueryRow(ctx, "select user_id, ip from sessions where refresh_token = $1 and expired_at > $2", hashed_token, time.Now()).Scan(&user_id, &prev_ip)
	if err != nil {
		trans.Rollback(ctx)
		return 0, "", err
	}
	trans.Commit(ctx)
	ans := uint32(user_id)
	return ans, prev_ip, nil
}

func (r *userRepo) CreateSession(ctx context.Context, userId uint32, refreshToken string, userIp string) error {
	trans, err := r.conn.Begin(ctx)
	if err != nil {
		trans.Rollback(ctx)
		return err
	}
	var success int
	h := sha256.New()
	h.Write([]byte(refreshToken))
	hashed_token := h.Sum(nil)
	err = trans.QueryRow(ctx, "insert into sessions (user_id, refresh_token, ip, expired_at) values($1, $2, $3,$4) on conflict (user_id) do update set refresh_token = $5, expired_at = $6 RETURNING 1;", userId, hashed_token, userIp, time.Now().Add(r.refreshTTL), refreshToken, time.Now().Add(r.refreshTTL)).Scan(&success)
	if err != nil {
		trans.Rollback(ctx)
		return err
	}
	if success != 1 {
		trans.Rollback(ctx)
		return fmt.Errorf("cant insert token")
	}
	trans.Commit(ctx)
	return nil
}

func (r *userRepo) CheckUserExist(ctx context.Context, id uint32) (bool, error) {
	trans, err := r.conn.Begin(ctx)
	if err != nil {
		trans.Rollback(ctx)
		return false, err
	}
	res := trans.QueryRow(ctx, "select id from users where id = $1", id)
	user_raw := &userSql{}
	res.Scan(&user_raw.Id)
	trans.Commit(ctx)
	if user_raw.Id.Valid {
		return true, nil
	}
	return false, ErrUserNotExist
}
