package domain

import "time"

type User struct {
	ID int

	Username     string `json,db:"username"`
	Email        string `json,db:"email"`
	PasswordHash string `json,db:"password_hash"`
	Phone        string `json,db:"phone"`

	Avatar string `json:"avatar"`

	FirstName  string `json,db:"first_name"`
	LastName   string `json,db:"last_name"`
	MiddleName string `json,db:"middle_name"`

	ConfirmToken string `json,db:"confirm_token"`
	IsConfirm    bool   `json,db:"is_confirm"`

	CreatedAt time.Time `json,db:"created_at"`
	UpdateAt  time.Time `json,db:"update_at"`

	Role UserRole `json,db:"role"`
}

type UserRole struct {
	ID   int    `json,db:"id"`
	Name string `json,db:"name"`
}
