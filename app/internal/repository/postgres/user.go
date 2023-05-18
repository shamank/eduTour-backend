package postgres

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/shamank/eduTour-backend/app/internal/domain"
)

type UserRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(ctx context.Context, user domain.User) error {

	q1 := `INSERT INTO USERS (username, email, password_hash) 
			VALUES ($1, $2, $3) RETURNING id`
	q2 := `INSERT INTO USERS_ROLES (user_id, role_id) VALUES ($1, 0)`
	tx, err := r.db.Begin()
	if err != nil {
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

func (r *UserRepo) GetByCredentials(ctx context.Context, email string, passwordHash string) (domain.User, error) {
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
		return domain.User{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var roles domain.UserRole
		rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&roles.ID,
			&roles.Name,
		)
		user.Roles = append(user.Roles, roles)
	}

	return user, nil
}

func (r *UserRepo) GetByRefreshToken(ctx context.Context, refreshToken string) (domain.User, error) {
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
	for rows.Next() {
		var userRole domain.UserRole
		rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&userRole.ID,
			&userRole.Name,
			&tokenID)
		user.Roles = append(user.Roles, userRole)
	}

	query2 := `UPDATE `
	//if err := r.db.Get(&user, query, email, passwordHash); err != nil {
	//	return domain.User{}, err
	//}

	return user, nil
}

func (r *UserRepo) SetRefreshToken(ctx context.Context, userID int, refreshInput domain.RefreshTokenInput) error {
	return nil
}

func (r *UserRepo) Verify(ctx context.Context, userID int) error {
	return nil
}
