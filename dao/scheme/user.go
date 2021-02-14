package scheme

import (
	"database/sql"
	uuid "github.com/jackc/pgtype/ext/gofrs-uuid"
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	ID              uint           `gorm:"primaryKey,autoIncrement" json:"id"`
	UUID            uuid.UUID      `gorm:"type:uuid,uniqueIndex" json:"uuid"`
	Name            string         `json:"name"`
	Slug            string         `json:"slug"`
	Password        string         `json:"-"`
	PasswordSalt    string         `json:"-"`
	Email           string         `json:"email"`
	Image           string         `json:"image"`
	Cover           string         `json:"cover"`
	BIO             string         `json:"bio"`
	Website         string         `json:"website"`
	Location        string         `json:"location"`
	Accessibility   string         `json:"accessibility"`
	Status          string         `json:"status"`
	Language        string         `json:"language"`
	MetaTitle       sql.NullString `json:"meta_title"`
	MetaDescription sql.NullString `json:"meta_description"`
	LastLogin       sql.NullTime   `json:"last_login"`
	CreatedAt       time.Time      `gorm:"autoCreateTime" json:"created_at"`
	CreatedBy       sql.NullInt64  `gorm:"not null" json:"created_by"`
	CreatedUser     *User          `gorm:"foreignKey:CreatedBy"`
	UpdatedAt       sql.NullTime   `gorm:"autoUpdateTime" json:"updated_at"`
	UpdatedBy       sql.NullInt64  `json:"updated_by"`
	UpdatedUser     *User          `gorm:"foreignKey:UpdatedBy"`
	Roles           []*Role        `gorm:"many2many:roles_users;"`
}

func (User) TableName() string {
	return "users"
}
