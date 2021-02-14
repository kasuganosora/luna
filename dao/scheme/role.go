package scheme

import (
	"database/sql"
	"github.com/google/uuid"
	"time"
)

type Role struct {
	ID          uint          `gorm:"primaryKey;autoIncrement" json:"id"`
	UUID        uuid.UUID     `gorm:"type:varchar(36);not null;uniqueIndex" json:"uuid"`
	Name        string        `gorm:"type:varchar(64);uniqueIndex" json:"name"`
	Description string        `gorm:"type:text" json:"description"`
	CreatedAt   time.Time     `gorm:"autoCreateTime" json:"created_at"`
	CreatedBy   sql.NullInt64 `gorm:"not null" json:"-"`
	CreatedUser *User         `gorm:"foreignKey:CreatedBy"`
	UpdatedAt   sql.NullTime  `gorm:"autoUpdateTime" json:"updated_at"`
	UpdatedBy   sql.NullInt64 `json:"-"`
	UpdatedUser *User         `gorm:"foreignKey:UpdatedBy"`
	Users       []*User       `gorm:"many2many:roles_users;"`
}

func (Role) TableName() string {
	return "roles"
}
