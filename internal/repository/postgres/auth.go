package postgres

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/shamank/eduTour-backend/internal/domain"
)

type AuthRepo struct {
	db *sqlx.DB
}

func NewAuthRepo(db *sqlx.DB) *AuthRepo {
	return &AuthRepo{db: db}
}

func (r *AuthRepo) Create(ctx context.Context, user domain.User) error {

	query := `INSERT INTO USERS (username, email, password_hash, confirm_token) 
			VALUES ($1, $2, $3, $4) RETURNING id`

	tx, err := r.db.Begin()
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			if pgErr.Code == "23505" {
				return domain.ErrUserAlreadyExists
			}
		}
		return err
	}
	var id int
	row := tx.QueryRow(query, user.Username, user.Email, user.PasswordHash, user.ConfirmToken)
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *AuthRepo) ConfirmUser(ctx context.Context, confirmToken string) error {

	query := `UPDATE USERS
				SET is_confirm = true
				WHERE confirm_token = $1`

	_, err := r.db.Exec(query, confirmToken)

	return err
}

func (r *AuthRepo) GetByCredentials(ctx context.Context, email string, passwordHash string) (domain.User, error) {
	var user domain.User

	query := `SELECT u.id, u.username, u.email, u.role_id, r.name
				FROM USERS u
				INNER JOIN ROLES r on u.role_id = r.id
				WHERE u.email = $1 AND u.password_hash = $2`

	err := r.db.QueryRow(query, email, passwordHash).Scan(&user.ID,
		&user.Username,
		&user.Email,
		&user.Role.ID,
		&user.Role.Name)

	return user, err
}

func (r *AuthRepo) GetByUsername(ctx context.Context, username string, passwordHash string) (domain.User, error) {
	var user domain.User

	query := `SELECT u.id, u.username, u.email, u.role_id, r.name
				FROM USERS u
				INNER JOIN ROLES r on u.role_id = r.id
				WHERE u.username = $1 AND u.password_hash = $2`

	err := r.db.QueryRow(query, username, passwordHash).Scan(&user.ID,
		&user.Username,
		&user.Email,
		&user.Role.ID,
		&user.Role.Name)

	return user, err
}

func (r *AuthRepo) GetByRefreshToken(ctx context.Context, refreshToken string) (domain.User, error) {
	var user domain.User
	var tokenID int
	query := `SELECT u.id, u.username, u.email, u.role_id, r.name, t.id
				FROM USERS u
				INNER JOIN ROLES r on r.id = u.role_id
				INNER JOIN REFRESH_TOKENS t on t.user_id = u.id
				WHERE t.refresh_token = $1 AND t.expire_at > CURRENT_TIMESTAMP AND NOT t.black_list`

	err := r.db.QueryRow(query, refreshToken).Scan(&user.ID,
		&user.Username,
		&user.Email,
		&user.Role.ID,
		&user.Role.Name,
		&tokenID)
	if err != nil {
		return domain.User{}, err
	}

	query2 := `UPDATE REFRESH_TOKENS
				SET black_list = true
				WHERE id = $1`
	_, err = r.db.Exec(query2, tokenID)
	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (r *AuthRepo) SetRefreshToken(ctx context.Context, userID int, refreshInput domain.RefreshTokenInput) error {
	query := `INSERT INTO REFRESH_TOKENS (user_id, refresh_token, expire_at) VALUES ($1, $2, to_timestamp($3))`

	_, err := r.db.Exec(query, userID, refreshInput.RefreshToken, int(refreshInput.ExpiresAt))

	return err
}

func (r *AuthRepo) Verify(ctx context.Context, userID int) error {
	return nil
}

func (r *AuthRepo) GetFullUserInfo(ctx context.Context, userID int) (domain.User, error) {

	var u domain.User

	query := `SELECT u.id, u.username, u.email, COALESCE(u.phone, '') as phone, 
				COALESCE(u.avatar, ''),  COALESCE(u.first_name, '') as first_name,
       COALESCE(u.last_name, '') as last_name, COALESCE(u.middle_name, '') as middle_name,
				u.created_at, r.id, r.name
				FROM USERS u
				INNER JOIN roles r on r.id = u.role_id
				WHERE u.id = $1`
	err := r.db.QueryRow(query, userID).Scan(&u.ID,
		&u.Username,
		&u.Email,
		&u.Phone,
		&u.Avatar,
		&u.FirstName,
		&u.LastName,
		&u.MiddleName,
		&u.CreatedAt,
		&u.Role.ID,
		&u.Role.Name)
	if err != nil {
		return domain.User{}, err
	}
	return u, nil
}
