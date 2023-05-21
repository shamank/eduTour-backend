package domain

import "time"

type User struct {
	ID int

	Username     string `json:"username"`
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
	Phone        string `json:"phone"`

	Avatar string `json:"avatar"`

	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	MiddleName string `json:"middle_name"`

	CreatedAt time.Time `json:"created_at"`
	//LastVisitAt time.Time

	Roles []UserRole `json:"roles"`
}

type UserRole struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
