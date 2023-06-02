package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/shamank/eduTour-backend/app/internal/domain"
)

type AuthRepo struct {
	db *sqlx.DB
}

func NewAuthRepo(db *sqlx.DB) *AuthRepo {
	return &AuthRepo{db: db}
}

func (r *AuthRepo) Create(ctx context.Context, user domain.User) error {

	q1 := `INSERT INTO USERS (username, email, password_hash) 
			VALUES ($1, $2, $3) RETURNING id`
	q2 := `INSERT INTO USERS_ROLES (user_id, role_id) VALUES ($1, 0)`
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
	row := tx.QueryRow(q1, user.Username, user.Email, user.PasswordHash)
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return err
	}
	if _, err := tx.Exec(q2, id); err != nil {
		return err
	}

	return tx.Commit()
}

func (r *AuthRepo) GetByCredentials(ctx context.Context, email string, passwordHash string) (domain.User, error) {
	var user domain.User

	query := `SELECT u.id, u.username, u.email, r.id, r.name
				FROM USERS u
				INNER JOIN USERS_ROLES ur on u.id = ur.user_id
				INNER JOIN ROLES r on r.id = ur.role_id
				WHERE u.email = $1 AND u.password_hash = $2`

	//if err := r.db.Get(&user, query, email, passwordHash); err != nil {
	//	return domain.User{}, err
	//}
	rows, err := r.db.Query(query, email, passwordHash)
	if err != nil {
		fmt.Printf("ERROR 1: %v", err)
		return domain.User{}, err
	}

	defer func() {
		if rErr := rows.Close(); rErr != nil {
			err = errors.Join(err, fmt.Errorf("error occured in closing row: %w", rErr))
		}
	}()

	found := false

	for rows.Next() {
		found = true
		var roles domain.UserRole
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&roles.ID,
			&roles.Name,
		)
		if err != nil {
			return domain.User{}, err
		}

		user.Roles = append(user.Roles, roles)
	}
	if !found {
		return domain.User{}, domain.ErrUserNotFound
	}

	return user, nil
}

func (r *AuthRepo) GetByRefreshToken(ctx context.Context, refreshToken string) (domain.User, error) {
	var user domain.User
	var tokenID int
	query := `SELECT u.id, u.username, u.email, r.id, r.name, t.id
				FROM USERS u
				INNER JOIN USERS_ROLES ur on u.id = ur.user_id
				INNER JOIN ROLES r on r.id = ur.role_id
				INNER JOIN REFRESH_TOKENS t on t.user_id = u.id
				WHERE t.refresh_token = $1 AND t.expire_at > CURRENT_TIMESTAMP AND NOT t.black_list`

	rows, err := r.db.Query(query, refreshToken)
	if err != nil {
		return domain.User{}, err
	}
	defer rows.Close()

	found := false
	for rows.Next() {
		found = true
		var userRole domain.UserRole
		if err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&userRole.ID,
			&userRole.Name,
			&tokenID); err != nil {
			return domain.User{}, err
		}
		user.Roles = append(user.Roles, userRole)
	}
	if !found {
		return domain.User{}, domain.ErrUserNotFound
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
				INNER JOIN users_roles ur on u.id = ur.user_id
				INNER JOIN roles r on r.id = ur.role_id
				WHERE u.id = $1`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return domain.User{}, err
	}
	defer rows.Close()

	found := false

	for rows.Next() {
		found = true
		var role domain.UserRole
		if err := rows.Scan(
			&u.ID,
			&u.Username,
			&u.Email,
			&u.Phone,
			&u.Avatar,
			&u.FirstName,
			&u.LastName,
			&u.MiddleName,
			&u.CreatedAt,
			&role.ID,
			&role.Name); err != nil {
			return domain.User{}, err
		}
		u.Roles = append(u.Roles, role)
	}

	if !found {
		return domain.User{}, domain.ErrUserNotFound
	}

	return u, nil
}
