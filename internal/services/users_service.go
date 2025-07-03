package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oThinas/bid/internal/store/pg"
	"golang.org/x/crypto/bcrypt"
)

const pgErrCodeUniqueViolation = "23505"

var ErrDuplicatedUsernameOrEmail = errors.New("username or email already exists")

type UserService struct {
	pool    *pgxpool.Pool
	queries *pg.Queries
}

func NewUserService(pool *pgxpool.Pool) UserService {
	return UserService{
		pool:    pool,
		queries: pg.New(pool),
	}
}

func (us *UserService) CreateUser(ctx context.Context, username, email, password, bio string) (uuid.UUID, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return uuid.Nil, err
	}

	args := pg.CreateUserParams{
		Username:     username,
		Email:        email,
		PasswordHash: hash,
		Bio:          bio,
	}

	id, err := us.queries.CreateUser(ctx, args)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgErrCodeUniqueViolation {
			return uuid.Nil, ErrDuplicatedUsernameOrEmail
		}

		return uuid.Nil, err
	}

	return id, nil
}
