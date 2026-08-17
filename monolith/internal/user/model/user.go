package model

import "time"

type User struct {
	ID           string    `db:"id"`
	FullName     string    `db:"username"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
	DeletedAt    time.Time `db:"deleted_at"`
}

type CreateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// type CreateUserRequest struct {
// 	Username string `json:"username" validate:"required,min=2,max=100"`
// 	Email    string `json:"email" validate:"required,email"`
// 	Password string `json:"password" validate:"required,min=8"`
// }

type UpdateUserRequest struct {
	FullName string `json:"full_name" binding:"required"`
}
