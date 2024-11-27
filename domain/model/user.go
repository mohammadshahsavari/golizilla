package model

import (
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User represents a user in the application.
type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"unique;not null" json:"username"`
	Email     string    `gorm:"unique;not null" json:"email"`
	Password  string    `gorm:"not null" json:"-"` // Exclude password from JSON responses
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// SetPassword hashes the password and assigns it to the Password field.
func (u *User) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword compares the provided password with the hashed password.
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// BeforeSave is a GORM hook that ensures the password is hashed before saving the record.
func (u *User) BeforeSave(tx *gorm.DB) error {
	fmt.Printf("BeforeSave Hook Triggered for User: %+v\n", u)
	if len(u.Password) > 0 && !isHashed(u.Password) {
		err := u.SetPassword(u.Password)
		if err != nil {
			fmt.Printf("Error in BeforeSave: %v\n", err)
			return err
		}
	}
	return nil
}

// isHashed checks if the password is already hashed.
func isHashed(password string) bool {
	return len(password) == 60 // bcrypt hashes are always 60 characters long
}

// ValidateEmail checks if the email field is valid.
func (u *User) ValidateEmail() error {
	if len(u.Email) == 0 {
		return errors.New("email is required")
	}
	// Add additional email validation logic if needed
	return nil
}
