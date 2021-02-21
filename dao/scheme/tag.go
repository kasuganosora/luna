package scheme

import (
	"database/sql"
	"github.com/google/uuid"
	"time"
)

type Tag struct {
	ID              uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	UUID            uuid.UUID      `gorm:"type:varchar(36);not null;uniqueIndex" json:"uuid"`
	Name            string         `gorm:"type:varchar(64);uniqueIndex" json:"name"`
	Slug            string         `gorm:"type:varchar(255)" json:"slug"`
	Description     string         `gorm:"type:text" json:"description"`
	ParentID        sql.NullInt64  `json:"parent_id"`
	MetaTitle       sql.NullString `gorm:"type:varchar(150);" json:"meta_title"`
	MetaDescription sql.NullString `gorm:"type:text" json:"meta_description"`
	CreatedAt       time.Time      `gorm:"autoCreateTime" json:"created_at"`
	CreatedBy       sql.NullInt64  `gorm:"not null" json:"-"`
	CreatedUser     *User          `gorm:"foreignKey:CreatedBy"`
	UpdatedAt       sql.NullTime   `gorm:"autoUpdateTime" json:"updated_at"`
	UpdatedBy       sql.NullInt64  `json:"updated_by" json:"-"`
	UpdatedUser     *User          `gorm:"foreignKey:UpdatedBy" json:"updated_user"`
	Posts           Posts          `gorm:"many2many:posts_tags" json:"posts"`
	Children        Tags           `gorm:"foreignkey:ParentID" json:"children"`
}

func (Tag) TableName() string {
	return "tags"
}

//go:generate pie Tags.*
type Tags []*Tag
