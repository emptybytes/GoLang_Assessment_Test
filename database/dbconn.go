// Package database provides functions for connecting to and interacting with the database.
// It also initializes the database connection and performs migrations.
package database

import (
	"fmt"
	"log"

	"github.com/alwilion/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Connection details for the PostgreSQL database
const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "admin123" // Enter your password for the DB
	dbname   = "product"
)

// Database source name (DSN) string for connecting to PostgreSQL
var dsn string = fmt.Sprintf("host=%s port=%d user=%s "+
	"password=%s dbname=%s sslmode=disable TimeZone=Asia/Shanghai",
	host, port, user, password, dbname)

// DB is a global variable representing the GORM database instance
var DB *gorm.DB

// DBconn initializes the database connection and performs migrations
func DBconn() {
	// Open a connection to the PostgreSQL database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	DB = db

	// Perform automatic migrations for the User and Product models
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Product{})
}
