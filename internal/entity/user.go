package entity

import "time"

type User struct {
	Name         string    `json:"name" validate:"required,max=50"`
	Surname      string    `json:"surname" validate:"required,max=50"`
	Contact      Contact   `json:"contact" validate:"required,dive,required"`
	PasswordHash string    `json:"passwordHash" validate:"required,max=255"`
	CreatedAt    time.Time `json:"createdAt" validate:"required"`
}

type Contact struct {
	Phone string `json:"phone" validate:"required,phone"`
	Email string `json:"email" validate:"required,email"`
}

//todo validate
