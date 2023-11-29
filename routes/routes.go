package routes

import (
	"github.com/alwilion/controllers"
	"github.com/gofiber/fiber/v2"
)

// Setup initializes the routes for the application
func Setup(app *fiber.App) {
	// Create a route group under the path "/user"
	api := app.Group("/user")

	// Define routes and associate them with corresponding controller functions
	api.Post("/login", controllers.Login)         // Route for user login
	api.Post("/register", controllers.Register)   // Route for user registration

	api.Get("/products", controllers.GetProductList)          // Route to get a list of products
	api.Get("/products/:id", controllers.GetProductById)      // Route to get a product by ID
	api.Delete("/products/:id", controllers.DeleteProductById) // Route to delete a product by ID
	api.Post("/products", controllers.AddProduct)              // Route to add a new product
	api.Put
