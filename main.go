package main

import (
	"fmt"

	"github.com/alwilion/database"
	"github.com/alwilion/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	// Print a message indicating the start of the application
	fmt.Println("Product Management")

	// Establish a connection to the database
	database.DBconn()

	// Create a new Fiber app instance
	app := fiber.New()

	// Use CORS middleware to handle Cross-Origin Resource Sharing
	app.Use(cors.New(cors.Config{
		AllowCredentials: true, // Important when using an HTTP-only cookie, allows frontend to access and send back the cookie
	}))

	// Set up routes for the application
	routes.Setup(app)

	// Start the application and listen on port 8000
	app.Listen(":8000")
}
