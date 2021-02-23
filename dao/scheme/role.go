package scheme

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Role struct {
	ID          uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	UUID        *uuid.UUID `gorm:"type:varchar(36);not null;uniqueIndex" json:"uuid"`
	Name        string     `gorm:"type:varchar(64);uniqueIndex" json:"name"`
	Description string     `gorm:"type:text" json:"description"`
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"created_at"`
	CreatedBy   *uint      `json:"-"`
	CreatedUser *User      `gorm:"foreignKey:CreatedBy"`
	UpdatedAt   *time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	UpdatedBy   *uint      `json:"-"`
	UpdatedUser *User      `gorm:"foreignKey:UpdatedBy"`
	Users       []*User    `gorm:"many2many:roles_users;"`
}

func (Role) TableName() string {
	return "roles"
}

func (r *Role) BeforeCreate(tx *gorm.DB) (err error) {
	if r.UUID == nil {
		roleUUID := uuid.New()
		r.UUID = &roleUUID
	}
	return
}
