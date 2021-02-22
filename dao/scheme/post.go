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
	UUID            *uuid.UUID     `gorm:"type:varchar(36);not null;uniqueIndex" json:"uuid"`
	Title           string         `gorm:"type:varchar(255); not null" json:"title"`
	Slug            *string        `gorm:"type:varchar(255);unique" json:"slug"`
	Markdown        string         `gorm:"type:longtext" json:"markdown"`
	HTML            string         `gorm:"type:longtext;column:html" json:"html"`
	Image           string         `gorm:"type:text" json:"image"`
	Featured        int            `gorm:"type:tinyint;not null;default:0" json:"featured"`
	Page            uint           `gorm:"type:tinyint;not null;default:0" json:"page"`
	Status          string         `gorm:"type:varchar(150);not null;default:draft; index" json:"status"`
	Language        string         `gorm:"type:varchar(20);not null;default:en_US'" json:"language"`
	MetaTitle       *string        `gorm:"type:varchar(150);" json:"meta_title"`
	MetaDescription *string        `gorm:"type:text" json:"meta_description"`
	AuthorID        *int64         `json:"author_id"`
	Author          *User          `gorm:"foreignKey:AuthorID"`
	CreatedAt       time.Time      `gorm:"autoCreateTime" json:"created_at"`
	CreatedBy       *int64         `json:"created_by"`
	CreatedUser     *User          `gorm:"foreignKey:CreatedBy"`
	UpdatedAt       *time.Time     `gorm:"autoUpdateTime" json:"updated_at"`
	UpdatedBy       sql.NullInt64  `json:"updated_by"`
	UpdatedUser     *User          `gorm:"foreignKey:UpdatedBy"`
	PublishedAt     *time.Time     `json:"published_at"`
	PublishedBy     *int64         `json:"published_by"`
	PublishedUser   *User          `gorm:"foreignKey:PublishedBy"`
	ScheduleTime    *time.Time     `json:"schedule_time"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Tags            Tags           `gorm:"many2many:posts_tags"`
}

func (Post) TableName() string {
	return "posts"
}

func (p *Post) FillFromMap(data map[string]interface{}) (err error) {
	return mapstructure.Decode(data, p)
}

func (p *Post) BeforeCreate(tx *gorm.DB) (err error) {
	if p.UUID == nil {
		postUUID := uuid.New()
		p.UUID = &postUUID
	}
	return
}

//go:generate pie Posts.*
type Posts []*Post
