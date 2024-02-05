package models

import (
	"gorm.io/gorm"
)

type Income struct {
	gorm.Model
	Amount string `json:"amount"` // string to avoid issues with it removing trailing zeros
	Source string `json:"source"`
}

type Expense struct {
	gorm.Model
	Amount string `json:"amount"`
	Source string `json:"source"`
}

type User struct {
	gorm.Model
	Username string `json:"username"`
	Password string `json:"password"`
	IsAdmin  bool   `json:"is_admin"`
}
