package entity

import "time"

type User struct {
	Id           int64     `json:"id" db:"id"`
	Name         string    `json:"name" validate:"required,max=50" db:"first_name"`
	Surname      string    `json:"surname" validate:"required,max=50" db:"last_name"`
	Contact      Contact   `json:"contact" validate:"required,dive,required"`
	PasswordHash string    `json:"-" validate:"required,max=255" db:"password_hash"`
	CreatedAt    time.Time `json:"createdAt" validate:"required" db:"created_at"`
	UpdatedAt    time.Time `json:"updatedAt" validate:"required" db:"updated_at"`
}

type Contact struct {
	Phone string `json:"phone" validate:"required,phone" db:"phone"`
	Email string `json:"email" validate:"required,email" db:"email"`
}

type SignUpInput struct {
	Password  string  `json:"password" validate:"required,min=8"`
	FirstName string  `json:"first_name" validate:"required"`
	LastName  string  `json:"last_name" validate:"required"`
	Contact   Contact `json:"contact" validate:"required,dive,required"`
}

type SignInInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

//todo validate
