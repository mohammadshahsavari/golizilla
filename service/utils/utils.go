package utils

import (
	"crypto/rand"
	"fmt"
)

func GenerateRandomCode(length int) string {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		// Handle error or panic
		panic(err)
	}
	// Convert bytes to digits
	code := ""
	for _, b := range bytes {
		code += fmt.Sprintf("%d", b%10)
	}
	return code
}
