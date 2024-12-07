package model

type Privilege struct {
	Id                    string `gorm:"primaryKey"`
	CanSetOnQuestionnaire bool
	Roles                 []*Role `gorm:"many2many:RolePrivilege"`
}
