package scheme

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"reflect"
	"strconv"
	"time"
)

type Setting struct {
	ID          uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	UUID        *uuid.UUID `gorm:"type:varchar(36);not null;uniqueIndex" json:"uuid"`
	Key         string     `gorm:"type:varchar(150); not null;uniqueIndex:type_key" json:"key"`
	Value       string     `gorm:"type:text" json:"value"`
	Type        string     `gorm:"type:varchar(150);default:core;not null;uniqueIndex:type_key" json:"type"`
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"created_at"`
	CreatedBy   *uint      `json:"-"`
	CreatedUser *User      `gorm:"foreignKey:CreatedBy" json:"created_user"`
	UpdatedAt   *time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	UpdatedBy   *uint      `json:"-"`
	UpdatedUser *User      `gorm:"foreignKey:UpdatedBy" json:"updated_user"`
}

func (Setting) TableName() string {
	return "settings"
}

func (s *Setting) BeforeCreate(tx *gorm.DB) (err error) {
	if s.UUID == nil {
		settingUUID := uuid.New()
		s.UUID = &settingUUID
	}
	return
}

func (s *Setting) SetValue(val interface{}) (err error) {
	switch val.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		v := reflect.ValueOf(val).Int()
		s.Value = strconv.FormatInt(v, 10)
	case float32, float64:
		v := reflect.ValueOf(val).Float()
		s.Value = strconv.FormatFloat(v, 'g', 32, 64)
	case bool:
		v := val.(bool)
		s.Value = strconv.FormatBool(v)
	case string:
		s.Value = val.(string)
	default:
		var b []byte
		b, err = json.Marshal(val)
		if err != nil {
			return
		}
		s.Value = string(b)
	}
	return
}

func (s *Setting) GetString() string {
	return s.Value
}

func (s *Setting) GetInt64() (i int64, err error) {
	i, err = strconv.ParseInt(s.Value, 10, 64)
	return
}

func (s *Setting) GetUInt() (i uint64, err error) {
	i, err = strconv.ParseUint(s.Value, 10, 64)
	return
}

func (s *Setting) GetFloat() (i float64, err error) {
	i, err = strconv.ParseFloat(s.Value, 10)
	return
}

func (s *Setting) GetBool() (i bool, err error) {
	i, err = strconv.ParseBool(s.Value)
	return
}

func (s *Setting) Scan(ptr interface{}) (err error) {
	err = json.Unmarshal([]byte(s.Value), ptr)
	return
}

func (s *Setting) CacheKey() string {
	return fmt.Sprintf("%s.%s", s.Type, s.Key)
}
