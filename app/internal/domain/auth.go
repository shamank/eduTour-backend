package domain

import "time"

type User struct {
	ID           int
	Username     string
	Email        string
	PasswordHash string
	CreatedAt    time.Time
	LastVisitAt  time.Time
	Roles        []UserRole
}

type UserRole struct {
	ID   int
	Name string
}
