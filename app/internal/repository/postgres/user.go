package postgres

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/shamank/eduTour-backend/app/internal/domain"
	"strings"
)

type UserRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (r *UserRepo) GetUserProfile(ctx context.Context, userName string) (domain.User, error) {

	var user domain.User

	query := `SELECT u.username, COALESCE(u.first_name, '') as first_name,
       COALESCE(u.last_name, '') as last_name, COALESCE(u.middle_name, '') as middle_name,
       COALESCE(u.avatar, '') as avatar, r.id, r.name
				FROM users u
				INNER JOIN users_roles ur on u.id = ur.user_id
				INNER JOIN roles r on r.id = ur.role_id
				WHERE u.username = $1`

	rows, err := r.db.Query(query, userName)
	if err != nil {
		return domain.User{}, err
	}
	defer rows.Close()

	found := false

	for rows.Next() {
		found = true
		var role domain.UserRole

		err := rows.Scan(&user.Username, &user.FirstName, &user.LastName, &user.MiddleName, &user.Avatar, &role.ID, &role.Name)
		if err != nil {
			return domain.User{}, err
		}
		user.Roles = append(user.Roles, role)
	}

	if !found {
		return domain.User{}, domain.ErrUserNotFound
	}

	return user, nil
}

func (r *UserRepo) UpdateUserProfile(ctx context.Context, user domain.User) error {

	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argID := 1

	if user.FirstName != "" {
		setValues = append(setValues, fmt.Sprintf("first_name=$%d", argID))
		args = append(args, user.FirstName)
		argID++
	}

	if user.LastName != "" {
		setValues = append(setValues, fmt.Sprintf("last_name=$%d", argID))
		args = append(args, user.LastName)
		argID++
	}

	if user.MiddleName != "" {
		setValues = append(setValues, fmt.Sprintf("middle_name=$%d", argID))
		args = append(args, user.MiddleName)
		argID++
	}

	if user.Avatar != "" {
		setValues = append(setValues, fmt.Sprintf("avatar=$%d", argID))
		args = append(args, user.Avatar)
		argID++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf(`UPDATE users SET %s WHERE username=$%d`, setQuery, argID)

	args = append(args, user.Username)

	_, err := r.db.Exec(query, args...)

	return err
}
