package scheme

import (
	"database/sql"
	"github.com/jackc/pgtype/ext/gofrs-uuid"
	"gorm.io/gorm"
	"time"
)

type Post struct {
	gorm.Model
	ID              uint           `gorm:"primaryKey,autoIncrement" json:"id"`
	UUID            uuid.UUID      `gorm:"type:uuid,uniqueIndex" json:"uuid"`
	Title           string         `json:"title"`
	Slug            sql.NullString `json:"slug"`
	Markdown        string         `gorm:"type:text" json:"markdown"`
	HTML            string         `gorm:"type:text,column:html" json:"html"`
	Image           string         `gorm:"type:text" json:"image"`
	Featured        int            `gorm:"type:tinyint not NULL,default:'0'" json:"featured"`
	Page            uint           `gorm:"type:tinyint not NULL,default:'0'" json:"page"`
	Status          string         `gorm:"type:varchar(150) not NULL,default:'draft'" json:"status"`
	Language        string         `gorm:"type:varchar(20) not NULL,default:'en_US'" json:"language"`
	MetaTitle       sql.NullString `json:"meta_title"`
	MetaDescription sql.NullString `json:"meta_description"`
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
	DeletedAt       gorm.DeletedAt `json:"deleted_at,omitempty"`
}

func (Post) TableName() string {
	return "posts"
}
