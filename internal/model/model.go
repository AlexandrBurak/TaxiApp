package model

// swagger:model Login
type Login struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// swagger:model User
type User struct {
	ID       string `json:"id" redis:"id"`
	Name     string `json:"name" redis:"name" binding:"required"`
	Phone    string `json:"phone" redis:"phone" binding:"required,e164"`
	Email    string `json:"email" redis:"email" binding:"required,email"`
	Password string `json:"password" redis:"password" binding:"required"`
}
