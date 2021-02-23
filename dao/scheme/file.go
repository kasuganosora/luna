package scheme

import (
	"github.com/google/uuid"
	"github.com/kabukky/journey/conversion"
	"gorm.io/gorm"
	"mime"
	"path/filepath"
	"time"
)

type File struct {
	ID           uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	UUID         *uuid.UUID `gorm:"type:varchar(36);not null;uniqueIndex" json:"uuid"`
	Name         string     `json:"name"`
	MIME         string     `json:"mime"`
	Size         int64      `json:"size"`
	Hash         string     `gorm:"type:varchar(255)" json:"hash"`
	Path         string     `json:"path"`
	AbsolutePath string     `json:"absolute_path"`
	CreatedBy    *uint      `json:"created_by"`
	CreatedUser  *User      `gorm:"foreignKey:CreatedBy"`
	CreatedAt    time.Time  `gorm:"autoCreateTime" json:"created_at"`
}

func (f *File) BeforeCreate(tx *gorm.DB) (err error) {
	if f.UUID == nil {
		postUUID := uuid.New()
		f.UUID = &postUUID
	}

	f.Name = conversion.XssFilter(f.Name)
	if f.MIME == "" {
		f.MIME = mime.TypeByExtension(filepath.Ext(f.Name))
	}

	return
}
