package utils

import (
	"crypto/rand"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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

func GenerateJWT(userID uuid.UUID, secretKey string, expiresIn time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"exp":     time.Now().Add(expiresIn).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseUUID(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}

// Input string in the format "yyyy-mm-dd"
func ParseDate(date string) (time.Time, error) {
	layout := "2006-01-02"
	
	parsedDate, err := time.Parse(layout, date)
	if err != nil {
		return time.Time{}, err
	}
	
	return parsedDate, nil
}