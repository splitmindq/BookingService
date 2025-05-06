package entity

type User struct {
	Name    string `json:"name" validate:"required,max=50"`
	Surname string `json:"surname" validate:"required,max=50"`
	Phone   string `json:"phone" validate:"required,phone"`
}

//todo validate
