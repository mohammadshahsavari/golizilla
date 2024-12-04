package model

import "github.com/google/uuid"

type RolePrivilegeOnInstance struct {
	Id              uuid.UUID `gorm:"primaryKey"`
	RoleId          uuid.UUID `gorm:"not null"`
	PrivilegeId     string
	QuestionnaireId uuid.UUID     `gorm:"not null"`
	Role            Role          `gorm:"foreinKey:RoleId"`
	Questionnaire   Questionnaire `gorm:"foreinKey:QuestionnaireId"`
}
