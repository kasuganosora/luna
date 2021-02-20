package setting

import (
	"github.com/kabukky/journey/structure"
	"time"
)

func insertSettingString(key string, value string, setting_type string, created_at time.Time, created_by int64) error {
}
func insertSettingInt64(key string, value int64, setting_type string, created_at time.Time, created_by int64) error {
}
func UpdateSettings(title []byte, description []byte, logo []byte, cover []byte, postsPerPage int64, activeTheme string, navigation []byte, updated_at time.Time, updated_by int64) error {
}
func RetrieveBlog() (*structure.Blog, error)                                             {}
func RetrieveActiveTheme() (*string, error)                                              {}
func UpdateActiveTheme(activeTheme string, updated_at time.Time, updated_by int64) error {}
