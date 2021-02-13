package scheme

import (
	"database/sql"
	uuid "github.com/jackc/pgtype/ext/gofrs-uuid"
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	ID              uint      `gorm:"primaryKey,autoIncrement" json:"id"`
	UUID            uuid.UUID `gorm:"type:uuid,uniqueIndex" json:"uuid"`
	Name            string
	Slug            string
	Password        string
	PasswordSalt    string
	Email           string
	Image           string
	Cover           string
	BIO             string
	Website         string
	Location        string
	Accessibility   string
	Status          string
	Language        string
	MetaTitle       sql.NullString `json:"meta_title"`
	MetaDescription sql.NullString `json:"meta_description"`
	LastLogin       sql.NullTime
	CreatedAt       time.Time     `gorm:"autoCreateTime" json:"created_at"`
	CreatedBy       sql.NullInt64 `gorm:"not null" json:"created_by"`
	CreatedUser     *User         `gorm:"foreignKey:CreatedBy"`
	UpdatedAt       sql.NullTime  `gorm:"autoUpdateTime" json:"updated_at"`
	UpdatedBy       sql.NullInt64 `json:"updated_by"`
	UpdatedUser     *User         `gorm:"foreignKey:UpdatedBy"`
}

func (User) TableName() string {
	return "users"
}
