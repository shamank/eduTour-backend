package domain

import "time"

type User struct {
	ID           int
	Username     string
	Email        string
	PasswordHash string
	CreatedAt    time.Time
	LastVisitAt  time.Time
}

type UserSignIn struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserSignUp struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type UserUpdate struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}
