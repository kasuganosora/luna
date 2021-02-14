package scheme

import (
	"database/sql"
	"github.com/google/uuid"
	"time"
)

type Setting struct {
	ID          uint          `gorm:"primaryKey;autoIncrement" json:"id"`
	UUID        uuid.UUID     `gorm:"type:varchar(36);not null;uniqueIndex" json:"uuid"`
	Key         string        `gorm:"type:varchar(150); not null" json:"key"`
	Value       string        `gorm:"type:text" json:"value"`
	Type        string        `gorm:"type:varchar(150);default:core;not null" json:"type"`
	CreatedAt   time.Time     `gorm:"autoCreateTime" json:"created_at"`
	CreatedBy   sql.NullInt64 `gorm:"not null" json:"-"`
	CreatedUser *User         `gorm:"foreignKey:CreatedBy" json:"created_user"`
	UpdatedAt   sql.NullTime  `gorm:"autoUpdateTime" json:"updated_at"`
	UpdatedBy   sql.NullInt64 `json:"-"`
	UpdatedUser *User         `gorm:"foreignKey:UpdatedBy" json:"updated_user"`
}

func (Setting) TableName() string {
	return "settings"
}
