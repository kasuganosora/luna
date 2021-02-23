package scheme

import (
	"crypto/rand"
	"encoding/hex"
	"github.com/google/uuid"
	"golang.org/x/crypto/scrypt"
	"gorm.io/gorm"
	"io"
	"strings"
	"time"
)

const PW_SALT_BYTES = 36
const PW_HASH_BYTES = 36

type User struct {
	ID              uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	UUID            *uuid.UUID     `gorm:"type:varchar(36);not null;uniqueIndex" json:"uuid"`
	Name            string         `gorm:"type:varchar(150)" json:"name"`
	Slug            string         `gorm:"type:varchar(255);unique" json:"slug"`
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
	MetaTitle       *string        `gorm:"type:varchar(150);" json:"meta_title"`
	MetaDescription *string        `gorm:"type:text" json:"meta_description"`
	LastLoginIP     *string        `json:"last_login_ip"`
	LastLoginTime   *time.Time     `json:"last_login_time"`
	CreatedAt       time.Time      `gorm:"autoCreateTime" json:"created_at"`
	CreatedBy       *uint          `json:"-"`
	CreatedUser     *User          `gorm:"foreignKey:CreatedBy"`
	UpdatedAt       *time.Time     `gorm:"autoUpdateTime" json:"updated_at"`
	UpdatedBy       *uint          `json:"-"`
	UpdatedUser     *User          `gorm:"foreignKey:UpdatedBy" json:"updated_user"`
	Roles           []*Role        `gorm:"many2many:roles_users" json:"roles"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

func (User) TableName() string {
	return "users"
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.UUID == nil {
		userUUID := uuid.New()
		u.UUID = &userUUID
	}
	return
}

func (u *User) SetPassword(password string) (err error) {
	password = strings.TrimSpace(password)
	salt, err := getSalt()
	if err != nil {
		return
	}

	hash, err := scrypt.Key([]byte(password), []byte(salt), 1<<15, 8, 1, PW_HASH_BYTES)
	if err != nil {
		return
	}
	u.Password = hex.EncodeToString(hash)
	u.PasswordSalt = salt
	return
}

func (u *User) ComparePassword(password string) (ok bool, err error) {
	password = strings.TrimSpace(password)
	hash, err := scrypt.Key([]byte(password), []byte(u.PasswordSalt), 1<<15, 8, 1, PW_HASH_BYTES)
	if err != nil {
		return
	}
	ok = hex.EncodeToString(hash) == u.Password
	return
}

func (u *User) UpdateLastLogin(db *gorm.DB, ip string, loginTime time.Time) (err error) {
	u.LastLoginIP = &ip
	u.LastLoginTime = &loginTime
	err = db.Save(&u).Error
	return
}

func getSalt() (salt string, err error) {
	saltByte := make([]byte, PW_SALT_BYTES)
	_, err = io.ReadFull(rand.Reader, saltByte)
	if err != nil {
		return
	}

	salt = hex.EncodeToString(saltByte)
	return
}
