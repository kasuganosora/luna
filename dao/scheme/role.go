package scheme

import (
	"database/sql"
	"github.com/gofrs/uuid"
	"time"
)

type Role struct {
	ID          uint          `gorm:"primaryKey,autoIncrement" json:"id"`
	UUID        uuid.UUID     `gorm:"type:uuid,uniqueIndex" json:"uuid"`
	Name        string        `gorm:"uniqueIndex" json:"name"`
	Description string        `json:"description"`
	CreatedAt   time.Time     `gorm:"autoCreateTime" json:"created_at"`
	CreatedBy   sql.NullInt64 `gorm:"not null" json:"-"`
	CreatedUser *User         `gorm:"foreignKey:CreatedBy"`
	UpdatedAt   sql.NullTime  `gorm:"autoUpdateTime" json:"updated_at"`
	UpdatedBy   sql.NullInt64 `json:"-"`
	UpdatedUser *User         `gorm:"foreignKey:UpdatedBy"`
	Users       []*User       `gorm:"many2many:roles_users;"`
}
