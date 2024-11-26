package persistence

import (
	"golizilla/domain/model"

	"gorm.io/gorm"
)

func MigrateRBAC(db *gorm.DB) error {
	err := db.AutoMigrate(&model.Role{}, &model.Permission{}, &model.UserRole{})
	if err != nil {
		return err
	}
	return nil
}
