package scheme

import (
	"database/sql"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
	"time"
)

type Tag struct {
	gorm.Model
	ID              uint           `gorm:"primaryKey,autoIncrement" json:"id"`
	UUID            uuid.UUID      `gorm:"type:uuid,uniqueIndex" json:"uuid"`
	Name            string         `gorm:"uniqueIndex" json:"name"`
	Slug            string         `json:"slug"`
	Description     string         `json:"description"`
	ParentID        sql.NullInt64  `json:"parent_id"`
	MetaTitle       sql.NullString `json:"meta_title"`
	MetaDescription sql.NullString `json:"meta_description"`
	CreatedAt       time.Time      `gorm:"autoCreateTime" json:"created_at"`
	CreatedBy       sql.NullInt64  `gorm:"not null" json:"-"`
	CreatedUser     *User          `gorm:"foreignKey:CreatedBy"`
	UpdatedAt       sql.NullTime   `gorm:"autoUpdateTime" json:"updated_at"`
	UpdatedBy       sql.NullInt64  `json:"updated_by" json:"-"`
	UpdatedUser     *User          `gorm:"foreignKey:UpdatedBy" json:"updated_user"`
	Posts           []*Post        `gorm:"many2many:posts_tags" json:"posts"`
}

func (Tag) TableName() string {
	return "tags"
}
