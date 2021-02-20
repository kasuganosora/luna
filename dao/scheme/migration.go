package scheme

import "time"

type Migration struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"type:varchar(255); not null;uniqueIndex"`
	Batch     int64     `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}
