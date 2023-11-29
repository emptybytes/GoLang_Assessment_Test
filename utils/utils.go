package utils

import (
	"math/rand"
	"time"
)

// GenerateRandomString generates a random string of the specified length
func GenerateRandomString(length int) string {
	// Seed the random number generator with the current Unix time
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)

	// Define the character set from which the random string will be generated
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// Create a byte slice of the specified length to store the random string
	result := make([]byte, length)

	// Populate each index in the byte slice with a randomly chosen character from the charset
	for i := range result {
		result[i] = charset[random.Intn(len(charset))]
	}

	// Convert the byte slice to a string and return the result
	return string(result)
}
