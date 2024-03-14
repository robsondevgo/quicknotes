package repositories

import (
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/robsondevgo/quicknotes/internal/models"
)

var ErrDuplicateEmail = newRepositoryError(errors.New("duplicate email"))
var ErrInvalidTokenOrUserAlreadyConfirmed = newRepositoryError(errors.New("invalid token or user already confirmed"))

type UserRepository interface {
	Create(ctx context.Context, email, password, token string) (*models.User, string, error)
	ConfirmUserByToken(ctx context.Context, token string) error
}

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{db: db}
}

func (ur *userRepository) Create(ctx context.Context, email, password, hashKey string) (*models.User, string, error) {
	var user models.User
	user.Email = pgtype.Text{String: email, Valid: true}
	user.Password = pgtype.Text{String: password, Valid: true}
	query := `INSERT INTO users (email, password)
			  VALUES ($1, $2)
			  RETURNING id, created_at`
	row := ur.db.QueryRow(ctx, query, user.Email, user.Password)
	if err := row.Scan(&user.Id, &user.CreatedAt); err != nil {
		if strings.Contains(err.Error(), "violates unique constraint") {
			return &user, "", ErrDuplicateEmail
		}
		return &user, "", newRepositoryError(err)
	}

	//gerar o token confirmation
	userToken, err := ur.createConfirmationToken(ctx, &user, hashKey)
	if err != nil {
		return nil, "", err
	}
	return &user, userToken.Token.String, nil
}

func (ur *userRepository) createConfirmationToken(ctx context.Context, user *models.User, token string) (*models.UserConfirmationToken, error) {
	var userToken models.UserConfirmationToken
	userToken.Token = pgtype.Text{String: token, Valid: true}
	userToken.UserId = user.Id
	query := `INSERT INTO users_confirmation_tokens (user_id, token)
			  VALUES ($1, $2)
			  RETURNING id, created_at`
	row := ur.db.QueryRow(ctx, query, userToken.UserId, userToken.Token)
	if err := row.Scan(&userToken.Id, &userToken.CreatedAt); err != nil {
		return nil, err
	}
	return &userToken, nil
}

func (ur *userRepository) ConfirmUserByToken(ctx context.Context, token string) error {
	//buscar o usuario e a confirmação do usuario pelo token
	query := `SELECT u.id u_id, t.id t_id FROM users u INNER JOIN users_confirmation_tokens t
			  ON u.id = t.user_id
			  WHERE u.active = false
			  AND t.confirmed = false
			  AND t.token = $1`
	row := ur.db.QueryRow(ctx, query, token)
	var userId, tokenId pgtype.Numeric
	err := row.Scan(&userId, &tokenId)
	if err != nil {
		if err == pgx.ErrNoRows {
			return ErrInvalidTokenOrUserAlreadyConfirmed
		}
		return newRepositoryError(err)
	}
	//tornar o usuário active
	queryUpdUser := "UPDATE users SET active = true, updated_at = now() WHERE id = $1"
	_, err = ur.db.Exec(ctx, queryUpdUser, userId)
	if err != nil {
		return newRepositoryError(err)
	}

	//tornar o token confirmed
	queryUpdToken := "UPDATE users_confirmation_tokens SET confirmed = true, updated_at = now() WHERE id = $1"
	_, err = ur.db.Exec(ctx, queryUpdToken, tokenId)
	if err != nil {
		return newRepositoryError(err)
	}

	return nil
}
