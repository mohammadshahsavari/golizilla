package model

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User represents a user in the domain layer.
type User struct {
	ID                      uuid.UUID `gorm:"type:uuid;primary_key;"`
	Username                string    `gorm:"unique;not null"`
	Email                   string    `gorm:"unique;not null"`
	NationalID              string    `gorm:"unique;not null"`
	Password                string    `gorm:"not null"`
	IsActive                bool      `gorm:"default:false"`
	EmailVerificationCode   string
	EmailVerificationExpiry time.Time
	IsTwoFAEnabled          bool `gorm:"default:false"`
	TwoFACode               string
	TwoFACodeExpiry         time.Time
	CreatedAt               time.Time
	UpdatedAt               time.Time
	FailedLoginAttempts     int  `gorm:"default:0"`
	AccountLocked           bool `gorm:"default:false"`
	AccountLockedUntil      time.Time
	// profile fields
	FirstName        string
	LastName         string
	City             string
	Wallet           uint
	DateOfBirth      time.Time       `gorm:"type:date"`
	NotificationList []*Notification `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
	RoleId           uuid.UUID       `gorm:"not null"`
	Role             Role            `gorm:"foreinKey:RoleId"`
}

// BeforeCreate is a GORM hook to generate a UUID before creating a new record.
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
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
