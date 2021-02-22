package scheme

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Tag struct {
	ID              uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	UUID            *uuid.UUID `gorm:"type:varchar(36);not null;uniqueIndex" json:"uuid"`
	Name            string     `gorm:"type:varchar(64);uniqueIndex" json:"name"`
	Slug            string     `gorm:"type:varchar(255)" json:"slug"`
	Description     string     `gorm:"type:text" json:"description"`
	ParentID        *uint      `json:"parent_id"`
	MetaTitle       *string    `gorm:"type:varchar(150);" json:"meta_title"`
	MetaDescription *string    `gorm:"type:text" json:"meta_description"`
	CreatedAt       time.Time  `gorm:"autoCreateTime" json:"created_at"`
	CreatedBy       *int64     `json:"-"`
	CreatedUser     *User      `gorm:"foreignKey:CreatedBy"`
	UpdatedAt       *time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	UpdatedBy       *int64     `json:"-"`
	UpdatedUser     *User      `gorm:"foreignKey:UpdatedBy" json:"updated_user"`
	Posts           Posts      `gorm:"many2many:posts_tags" json:"posts"`
	Children        Tags       `gorm:"foreignkey:ParentID" json:"children"`
}

func (Tag) TableName() string {
	return "tags"
}

func (t *Tag) BeforeCreate(tx *gorm.DB) (err error) {
	if t.UUID == nil {
		tagUUID := uuid.New()
		t.UUID = &tagUUID
	}
	return
}

//go:generate pie Tags.*
type Tags []*Tag
