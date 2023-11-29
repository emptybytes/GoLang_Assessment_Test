package controllers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/alwilion/database"
	"github.com/alwilion/models"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

// SecretKey is the secret key used for JWT token creation and validation
const SecretKey = "secret"

// Register handles user registration
func Register(c *fiber.Ctx) error {
	var data map[string]string

	// Parse request body into a map
	if err := c.BodyParser(&data); err != nil {
		return err
	}

	// Check if password is provided
	if data["password"] == "" {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Invalid password",
		})
	}

	// Hash the password
	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)

	// Create a new user with the provided data
	user := models.User{
		Name:     data["name"],
		Email:    data["email"],
		Password: password,
	}

	// Validate the user struct using the validator package
	validate := validator.New()
	err := validate.Struct(user)

	// Handle validation errors
	if err != nil {
		errorMsg := ""
		for _, err := range err.(validator.ValidationErrors) {
			fmt.Println("field:", err.Field())
			if err.Field() == "Email" {
				errorMsg = "Issue in Email"
			} else if err.Field() == "Password" {
				errorMsg = "Issue in Password"
			} else if err.Field() == "Name" {
				errorMsg = "Issue in Name Field"
			}
		}
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": errorMsg,
		})
	}

	// Insert the user into the database
	result := database.DB.Create(&user)

	// Check for errors during insertion
	if result.Error != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "failed to insert record into the database",
		})
	}
	return c.JSON(&user)
}

// Login handles user login
func Login(c *fiber.Ctx) error {
	var data map[string]string

	// Parse request body into a map
	if err := c.BodyParser(&data); err != nil {
		return err
	}

	// Find the user in the database by email
	var user models.User
	database.DB.Where("email = ?", data["email"]).First(&user)

	// Check if the user exists
	if user.ID == 0 {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{"message": "User not found"})
	}

	// Compare the provided password with the hashed password in the database
	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"])); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "incorrect password",
		})
	}

	// Create a new JWT token with user ID as the issuer and an expiration time of 24 hours
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(user.ID)),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	})

	// Sign the token with the secret key
	token, err := claims.SignedString([]byte(SecretKey))

	// Check for errors during token creation
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "could not login",
		})
	}

	// Return the success message along with the generated token
	return c.JSON(fiber.Map{
		"message": "success",
		"token":   token,
	})
}

// User retrieves user details based on the provided JWT token
func User(c *fiber.Ctx) error {
	// Authenticate the request and retrieve the JWT token
	token, err := authentication(c)

	// Handle authentication errors
	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}

	// Extract the claims from the JWT token
	claims := token.Claims.(*jwt.StandardClaims)

	// Retrieve the user from the database based on the user ID in the JWT claims
	var user models.User
	database.DB.Where("id = ?", claims.Issuer).First(&user)

	// Print the type of user ID for debugging
	fmt.Printf("type of id %T", user.ID)

	// Return the user details as JSON
	return c.JSON(user)
}

// AddProduct handles the addition of a new product
func AddProduct(c *fiber.Ctx) error {
	fmt.Print("Add Product")

	// Authenticate the request
	_, err := authentication(c)

	// Handle authentication errors
	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}

	// Parse the request body into a map
	data := make(map[string]interface{})
	if err := c.BodyParser(&data); err != nil {
		return err
	}

	fmt.Print("Add Product1")

	// Create a new product instance
	var product models.Product

	// Validate and set the product fields from the parsed data
	if value, ok := data["price"]; ok {
		// Type assertion to float64
		price, ok := value.(float64)
		if !ok || price <= 0 {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"message": "Invalid or missing positive Price",
			})
		}
		product.Price = price
	} else {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Price Key Missing",
		})
	}

	if value, ok := data["description"]; ok {
		// Type assertion to string
		description, ok := value.(string)
		if !ok {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"message": "Invalid Description Key type",
			})
		}
		product.Description = description
	} else {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Description Key Missing",
		})
	}

	if value, ok := data["name"]; ok {
		// Type assertion to string
		name, ok := value.(string)
		if !ok {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"message": "Invalid Name Key type",
			})
		}
		product.Name = name
	} else {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Name Key Missing",
		})
	}

	// Insert the product into the database
	result := database.DB.Create(&product)

	// Check for errors during insertion
	if result.Error != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "failed to insert record into the database",
		})
	}
	return c.JSON(product)
}

// UpdateProduct handles the update of an existing product
func UpdateProduct(c *fiber.Ctx) error {
	// Authenticate the request
	_, err := authentication(c)

	// Handle authentication errors
	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}

	var updatedProduct models.Product

	if err := database.DB.First(&updatedProduct, c.Params("id")).Error; err != nil {
		return c.JSON(fiber.Map{"message": "Product Not Found"})
	}

	data := make(map[string]interface{})
	if err := c.BodyParser(&data); err != nil {
		return err
	}

	// Validate and update the product fields from the parsed data
	if value, ok := data["price"]; ok {
		// Type assertion to float64
		price, ok := value.(float64)
		if !ok || price <= 0 {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"message": "Invalid or missing positive Price",
			})
		}
		updatedProduct.Price = price
	}
	if value, ok := data["description"]; ok {
		// Type assertion to string
		description, ok := value.(string)
		if !ok {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"message": "Invalid Description Key type",
			})
		}
		updatedProduct.Description = description
	}
	if value, ok := data["name"]; ok {
		// Type assertion to string
		name, ok := value.(string)
		if !ok {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"message": "Invalid Name Key type",
			})
		}
		updatedProduct.Name = name
	}

	// Save the updated product in the database
	result := database.DB.Save(&updatedProduct)

	// Check for errors during update
	if result.Error != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "failed to update record into the database",
		})
	}
	return c.JSON(updatedProduct)
}
