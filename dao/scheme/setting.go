package scheme

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Setting struct {
	ID          uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	UUID        *uuid.UUID `gorm:"type:varchar(36);not null;uniqueIndex" json:"uuid"`
	Key         string     `gorm:"type:varchar(150); not null;uniqueIndex:type_key" json:"key"`
	Value       string     `gorm:"type:text" json:"value"`
	Type        string     `gorm:"type:varchar(150);default:core;not null;uniqueIndex:type_key" json:"type"`
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"created_at"`
	CreatedBy   *int64     `json:"-"`
	CreatedUser *User      `gorm:"foreignKey:CreatedBy" json:"created_user"`
	UpdatedAt   *time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	UpdatedBy   *int64     `json:"-"`
	UpdatedUser *User      `gorm:"foreignKey:UpdatedBy" json:"updated_user"`
}

func (Setting) TableName() string {
	return "settings"
}

func (s *Setting) BeforeCreate(tx *gorm.DB) (err error) {
	if s.UUID == nil {
		settingUUID := uuid.New()
		s.UUID = &settingUUID
	}
	return
}
