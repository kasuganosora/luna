package scheme

import (
	"database/sql"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID              uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	UUID            uuid.UUID      `gorm:"type:varchar(36);not null;uniqueIndex" json:"uuid"`
	Name            string         `gorm:"type:varchar(150)" json:"name"`
	Slug            string         `gorm:"type:varchar(255)" json:"slug"`
	Password        string         `gorm:"type:varchar(255)" json:"-"`
	PasswordSalt    string         `gorm:"type:varchar(255)" json:"-"`
	Email           string         `gorm:"type:varchar(255)" json:"email"`
	Image           string         `gorm:"type:text" json:"image"`
	Cover           string         `gorm:"type:text" json:"cover"`
	BIO             string         `gorm:"type:varchar(255)" json:"bio"`
	Website         string         `gorm:"type:varchar(255)" json:"website"`
	Location        string         `gorm:"type:text" json:"location"`
	Accessibility   string         `gorm:"type:text" json:"accessibility"`
	Status          string         `gorm:"type:varchar(160); not null; default:active" json:"status"`
	Language        string         `gorm:"type:varchar(6); not null; default:en_US" json:"language"`
	MetaTitle       sql.NullString `gorm:"type:varchar(150);" json:"meta_title"`
	MetaDescription sql.NullString `gorm:"type:text" json:"meta_description"`
	LastLogin       sql.NullTime   `json:"last_login"`
	CreatedAt       time.Time      `gorm:"autoCreateTime" json:"created_at"`
	CreatedBy       sql.NullInt64  `json:"-"`
	CreatedUser     *User          `gorm:"foreignKey:CreatedBy"`
	UpdatedAt       sql.NullTime   `gorm:"autoUpdateTime" json:"updated_at"`
	UpdatedBy       sql.NullInt64  `json:"-"`
	UpdatedUser     *User          `gorm:"foreignKey:UpdatedBy" json:"updated_user"`
	Roles           []*Role        `gorm:"many2many:roles_users" json:"roles"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

func (User) TableName() string {
	return "users"
}
