package setting

import (
	"github.com/kabukky/journey/dao/scheme"
	"github.com/kabukky/journey/repositories/post"
	"github.com/kabukky/journey/structure"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"sync"
)

/*
func insertSettingString(key string, value string, setting_type string, created_at time.Time, created_by int64) error {
}
func insertSettingInt64(key string, value int64, setting_type string, created_at time.Time, created_by int64) error {
}
func UpdateSettings(title []byte, description []byte, logo []byte, cover []byte, postsPerPage int64, activeTheme string, navigation []byte, updated_at time.Time, updated_by int64) error {
}
func RetrieveBlog() (*structure.Blog, error)                                             {}
func RetrieveActiveTheme() (*string, error)                                              {}
func UpdateActiveTheme(activeTheme string, updated_at time.Time, updated_by int64) error {}
*/

var settingCache = sync.Map{}

func Set(db gorm.DB, settingType string, key string, value interface{}, user *scheme.User) (err error) {
	var setting *scheme.Setting
	var isNew bool
	err = db.Model(&scheme.Setting{}).
		Where("type = ?", settingType).
		Where("name = ?", key).
		First(&setting).Error

	if err != nil && !errors.Is(gorm.ErrRecordNotFound, err) {
		return
	}

	if setting == nil {
		setting = &scheme.Setting{}
		isNew = true
	}
	setting.Key = key
	setting.Type = settingType
	if err = setting.SetValue(value); err != nil {
		return
	}

	if user != nil {
		if isNew {
			setting.CreatedBy = &user.ID
		} else {
			setting.UpdatedBy = &user.ID
		}
	}

	if err = db.Save(setting).Error; err != nil {
		return
	}

	settingCache.Store(setting.CacheKey(), setting)
	return
}

func Get(db *gorm.DB, settingType string, key string) (setting *scheme.Setting, err error) {
	setting = &scheme.Setting{}
	setting.Key = key
	setting.Type = settingType
	tmp, ok := settingCache.Load(setting.CacheKey())
	if ok {
		setting = tmp.(*scheme.Setting)
		return
	}

	err = db.Model(&scheme.Setting{}).
		Where("type = ?", settingType).
		Where("name = ?", key).
		First(&setting).Error

	if err != nil {
		return
	}

	settingCache.Store(setting.CacheKey(), setting)
	return
}

func LoadCache(db *gorm.DB) (err error) {
	settingCache = sync.Map{}
	settings := make([]*scheme.Setting, 0)
	err = db.Model(&scheme.Setting{}).Find(&settings).Error
	for _, setting := range settings {
		settingCache.Store(setting.CacheKey(), setting)
	}
	return
}

func RetrieveBlog(db *gorm.DB) (blog *structure.Blog, err error) {
	blog = &structure.Blog{}
	//Url             string
	//Title           string
	//Description     string
	//Logo            string
	//Cover           string
	//AssetPath       string
	//PostCount       int64
	//PostsPerPage    int64
	//ActiveTheme     string
	//NavigationItems []Navigation

	var setting *scheme.Setting
	if setting, err = Get(db, "blog", "title"); err != nil {
		return
	}
	blog.Title = setting.GetString()

	if setting, err = Get(db, "blog", "description"); err != nil {
		return
	}
	blog.Description = setting.GetString()

	if setting, err = Get(db, "blog", "logo"); err != nil {
		return
	}
	blog.Logo = setting.GetString()

	if setting, err = Get(db, "blog", "cover"); err != nil {
		return
	}
	blog.Cover = setting.GetString()

	if setting, err = Get(db, "blog", "postsPerPage"); err != nil {
		return
	}
	blog.PostsPerPage, err = setting.GetInt64()
	if err != nil {
		return
	}

	if setting, err = Get(db, "blog", "activeTheme"); err != nil {
		return
	}
	blog.ActiveTheme = setting.GetString()

	blog.NavigationItems = make([]structure.Navigation, 0)
	if setting, err = Get(db, "blog", "navigation"); err != nil {
		return
	}

	if err = setting.Scan(&(blog.NavigationItems)); err != nil {
		return
	}

	blog.PostCount, err = post.GetTotalPostCount(db)
	if err != nil {
		return
	}

	return
}
