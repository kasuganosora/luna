package scheme

import (
	"github.com/google/uuid"
	"github.com/kabukky/journey/conversion"
	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
	"time"
)

const POST_STATUS_DRAFT = "draft"
const POST_STATUS_PUBLISHED = "published"

type Post struct {
	ID              uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	UUID            *uuid.UUID     `gorm:"type:varchar(36);not null;uniqueIndex" json:"uuid"`
	Title           string         `gorm:"type:varchar(255); not null" json:"title"`
	Slug            *string        `gorm:"type:varchar(255);unique" json:"slug"`
	Markdown        string         `gorm:"type:longtext" json:"markdown"`
	HTML            string         `gorm:"type:longtext;column:html" json:"html"`
	Image           string         `gorm:"type:text" json:"image"`
	Featured        int            `gorm:"type:tinyint;not null;default:0" json:"featured"`
	Page            bool           `gorm:"type:tinyint;not null;default:0" json:"page"`
	Status          string         `gorm:"type:varchar(150);not null;default:draft; index" json:"status"`
	Language        string         `gorm:"type:varchar(20);not null;default:en_US'" json:"language"`
	MetaTitle       *string        `gorm:"type:varchar(150);" json:"meta_title"`
	MetaDescription *string        `gorm:"type:text" json:"meta_description"`
	AuthorID        *uint          `json:"author_id"`
	Author          *User          `gorm:"foreignKey:AuthorID" json:"author"`
	CreatedAt       time.Time      `gorm:"autoCreateTime" json:"created_at"`
	CreatedBy       *uint          `json:"created_by"`
	CreatedUser     *User          `gorm:"foreignKey:CreatedBy"`
	UpdatedAt       *time.Time     `gorm:"autoUpdateTime" json:"updated_at"`
	UpdatedBy       *uint          `json:"updated_by"`
	UpdatedUser     *User          `gorm:"foreignKey:UpdatedBy"`
	PublishedAt     *time.Time     `json:"published_at"`
	PublishedBy     *uint          `json:"published_by"`
	PublishedUser   *User          `gorm:"foreignKey:PublishedBy"`
	ScheduleTime    *time.Time     `json:"schedule_time"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Tags            Tags           `gorm:"many2many:posts_tags" json:"tags"`
	TagsStr         string         `gorm:"-" json:"tags_str,omitempty"`
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

	p.saveEvent(tx)

	return
}

func (p *Post) BeforeUpdate(tx *gorm.DB) (err error) {
	p.saveEvent(tx)
	return
}

func (p *Post) IsPage() bool {
	return p.Page
}

func (p *Post) safeStringClean() {
	p.Title = conversion.XssFilter(p.Title)
	p.HTML = conversion.XssFilter(p.HTML)
	if p.MetaTitle != nil {
		mt := conversion.XssFilter(*p.MetaTitle)
		p.MetaTitle = &mt
	}
	if p.MetaDescription != nil {
		md := conversion.XssFilter(*p.MetaDescription)
		p.MetaDescription = &md
	}

	if p.Slug != nil {
		slug := conversion.XssFilter(*p.Slug)
		p.Slug = &slug
	}
}

func (p *Post) saveEvent(tx *gorm.DB) {
	p.safeStringClean()

	if p.Status == POST_STATUS_PUBLISHED && p.PublishedAt == nil {
		now := time.Now()
		p.PublishedAt = &now
	}

	if p.Status == POST_STATUS_DRAFT {
		p.PublishedAt = nil
	}
}

//go:generate pie Posts.*
type Posts []*Post
