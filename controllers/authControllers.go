package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

// authentication is a helper function to parse and validate JWT tokens from the Authorization header
func authentication(c *fiber.Ctx) (*jwt.Token, error) {
	// Get the Authorization header value, which should contain the JWT token
	authorizationHeader := c.Get("Authorization")

	// Parse and validate the JWT token with the specified claims and key
	token, err := jwt.ParseWithClaims(authorizationHeader, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil // Using the SecretKey which was generated in the Login function
	})

	// Return the parsed token and any parsing error
	return token, err
}
