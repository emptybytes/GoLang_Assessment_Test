package models

import "gorm.io/gorm"

// User represents the model for user data
type User struct {
	gorm.Model
	Name     string `json:"name" validate:"required"`                   // User's name
	Email    string `json:"email" gorm:"unique" validate:"required,email"` // User's email (unique constraint)
	Password []byte `json:"password" validate:"required"`                // User's hashed password
}

// Product represents the model for product data
type Product struct {
	gorm.Model
	Name        string  `json:"name" validate:"required"`          // Product name
	Description string  `json:"description" validate:"required"`   // Product description
	Price       float64 `json:"price" validate:"required"`         // Product price
}
