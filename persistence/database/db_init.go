package database

import (
	"errors"
	"fmt"
	"golizilla/config"
	models "golizilla/domain/model"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// SetupDB initializes and returns a database connection
func setupDB(cfg *config.Config) (*gorm.DB, error) {
	// Retrieve the database URL from environment variables
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", cfg.DBUsername, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	// Connect to the database using GORM
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}

	// Run migrations for the Ads, Filters, and Users models
	if err := db.AutoMigrate(&models.Question{}); err != nil {
		log.Fatalf("Failed to migrate Question: %v", err)
	}
	if err := db.AutoMigrate(&models.Option{}); err != nil {
		log.Fatalf("Failed to migrate Option: %v", err)
	}
	err = db.AutoMigrate(
		&models.User{},
		&models.Questionnaire{},
		&models.Notification{},
		&models.Question{},
		&models.Option{},
		&models.Answer{},
		&models.Role{},
		&models.Privilege{},
		&models.RolePrivilege{},
		&models.RolePrivilegeOnInstance{},
		&models.UserSubmission{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Check if Super Admin exists, and create one if not
	err = db.Where("username = ?", cfg.AdminUsername).First(&models.User{}).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			createSuperAdmin(db, cfg)
		} else {
			return nil, err
		}
	}

	err = createPrivileges(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func SetupDatabaseWithRetry(cfg *config.Config, retries int) (interface{}, error) {
	var db interface{}
	var err error
	for retries > 0 {
		db, err = setupDB(cfg)
		if err == nil {
			return db, nil
		}
		retries--
		log.Printf("Error setting up the database, retrying... (%d retries left)", retries)
		time.Sleep(2 * time.Second)
	}
	return nil, err
}

func createSuperAdmin(db *gorm.DB, cfg *config.Config) error {
	role := &models.Role{
		Name: "super-admin",
	}
	if err := db.Create(role).Error; err != nil {
		return err
	}
	admin := &models.User{
		Username:   cfg.AdminUsername,
		Email:      cfg.AdminEmail,
		NationalID: cfg.AdminNationalID,
		Password:   cfg.AdminPassword,
		RoleId:     role.ID,
	}
	return db.Create(admin).Error
}

func createPrivileges(db *gorm.DB) error {
	// Define a list of privileges
	privileges := []models.Privilege{
		{Id: "CreateQuestionnaire"},
		{Id: "EditQuestionnaire"},
		{Id: "DeleteQuestionnaire"},
		{Id: "ViewQuestionnaire"},
		{Id: "CreateQuestion"},
		{Id: "EditQuestion"},
		{Id: "DeleteQuestion"},
		{Id: "ViewQuestion"},
		{Id: "AssignRole"},
		{Id: "ManagePrivileges"},
	}

	// Loop through the privileges and add them if they do not exist
	for _, privilege := range privileges {
		var existingPrivilege models.Privilege
		err := db.Where("id = ?", privilege.Id).First(&existingPrivilege).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Privilege does not exist, create it
			if err := db.Create(&privilege).Error; err != nil {
				log.Printf("Failed to create privilege '%s': %v", privilege.Id, err)
				return err
			}
			log.Printf("Privilege '%s' created successfully.", privilege.Id)
		} else if err != nil {
			// An error occurred while querying the database
			log.Printf("Error checking privilege '%s': %v", privilege.Id, err)
			return err
		}
	}

	return nil
}
