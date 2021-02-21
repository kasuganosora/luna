package scheme

import (
	"database/sql"
	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
	"time"
)

type Post struct {
	ID              uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	UUID            uuid.UUID      `gorm:"type:varchar(36);not null;uniqueIndex" json:"uuid"`
	Title           string         `gorm:"type:varchar(255); not null" json:"title"`
	Slug            sql.NullString `gorm:"type:varchar(255)" json:"slug"`
	Markdown        string         `gorm:"type:longtext" json:"markdown"`
	HTML            string         `gorm:"type:longtext;column:html" json:"html"`
	Image           string         `gorm:"type:text" json:"image"`
	Featured        int            `gorm:"type:tinyint;not null;default:0" json:"featured"`
	Page            uint           `gorm:"type:tinyint;not null;default:0" json:"page"`
	Status          string         `gorm:"type:varchar(150);not null;default:draft; index" json:"status"`
	Language        string         `gorm:"type:varchar(20);not null;default:en_US'" json:"language"`
	MetaTitle       sql.NullString `gorm:"type:varchar(150);" json:"meta_title"`
	MetaDescription sql.NullString `gorm:"type:text" json:"meta_description"`
	AuthorID        int            `gorm:"not null" json:"author_id"`
	Author          *User          `gorm:"foreignKey:AuthorID"`
	CreatedAt       time.Time      `gorm:"autoCreateTime" json:"created_at"`
	CreatedBy       sql.NullInt64  `gorm:"not null" json:"created_by"`
	CreatedUser     *User          `gorm:"foreignKey:CreatedBy"`
	UpdatedAt       sql.NullTime   `gorm:"autoUpdateTime" json:"updated_at"`
	UpdatedBy       sql.NullInt64  `json:"updated_by"`
	UpdatedUser     *User          `gorm:"foreignKey:UpdatedBy"`
	PublishedAt     sql.NullTime   `json:"published_at"`
	PublishedBy     sql.NullInt64  `json:"published_by"`
	PublishedUser   *User          `gorm:"foreignKey:PublishedBy"`
	ScheduleTime    *sql.NullTime  `json:"schedule_time"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Tags            Tags           `gorm:"many2many:posts_tags"`
}

func (Post) TableName() string {
	return "posts"
}

func (p *Post) FillFromMap(data map[string]interface{}) (err error) {
	return mapstructure.Decode(data, p)
}

//go:generate pie Posts.*
type Posts []*Post
