package setting

import (
	"github.com/joho/godotenv"
	"github.com/kabukky/journey/dao/scheme"
	"github.com/kabukky/journey/flags"
	"github.com/kabukky/journey/repositories/post"
	"github.com/kabukky/journey/structure"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var settingCache = sync.Map{}
var ErrKeyNotExists = errors.New("key is not exists")

func Set(db *gorm.DB, settingType string, key string, value interface{}, user *scheme.User) (err error) {
	var setting *scheme.Setting

	// global setting set is temporary set to cache, but isn't real save to disk.
	if settingType == "global" {
		setting = &scheme.Setting{}
		setting.Type = "global"
		setting.Key = strings.TrimSpace(key)
		err = setting.SetValue(value)
		if err != nil {
			return
		}
		settingCache.Store(setting.CacheKey(), setting)
		return
	}

	//

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

func GetGlobal(key string) (setting *scheme.Setting, err error) {
	setting = &scheme.Setting{}
	setting.Key = key
	setting.Type = "global"

	tmp, ok := settingCache.Load(setting.CacheKey())
	if ok {
		setting = tmp.(*scheme.Setting)
		return
	}
	err = ErrKeyNotExists
	return
}

func LoadCache(db *gorm.DB) (err error) {
	settingCache = sync.Map{}
	settings := make([]*scheme.Setting, 0)
	err = db.Model(&scheme.Setting{}).Find(&settings).Error
	for _, setting := range settings {
		settingCache.Store(setting.CacheKey(), setting)
	}

	LoadEnv()
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

	blog.AssetPath = "/assets/"

	return
}

func PostsPerPage(db *gorm.DB, defaultVal int64) (val int64, err error) {
	postsPerPageSetting, err := Get(db, "blog", "postsPerPage")
	if err != nil && !errors.Is(gorm.ErrRecordNotFound, err) {
		return
	}
	if postsPerPageSetting == nil {
		val = defaultVal
		return
	}

	val, err = postsPerPageSetting.GetInt64()
	if err != nil {
		val = defaultVal
		err = nil
		return
	}
	return
}

func RetrieveActiveTheme(db *gorm.DB) (setting *scheme.Setting, err error) {
	return Get(db, "blog", "activeTheme")
}

func UpdateActiveTheme(db *gorm.DB, themeName string, user *scheme.User) (err error) {
	return Set(db, "blog", "activeTheme", themeName, user)
}

func LoadEnv() {
	envs := os.Environ()
	for _, ev := range envs {
		kv := strings.Split(ev, "=")
		if len(kv) != 2 {
			continue
		}
		setting := &scheme.Setting{}
		setting.Type = "global"
		setting.Key = strings.TrimSpace(kv[0])
		setting.Value = strings.TrimSpace(kv[1])
		settingCache.Store(setting.CacheKey(), setting)
	}

	envFile, err := getEnvFile()
	if err == nil {
		var envContent map[string]string
		envContent, err = godotenv.Read(envFile)
		if err != nil {
			return
		}
		for k, v := range envContent {
			setting := &scheme.Setting{}
			setting.Type = "global"
			setting.Key = strings.TrimSpace(k)
			setting.Value = strings.TrimSpace(v)
			settingCache.Store(setting.CacheKey(), setting)
		}
	}

	for k, v := range flags.Settings {
		setting := &scheme.Setting{}
		setting.Type = "global"
		setting.Key = strings.TrimSpace(k)
		_ = setting.SetValue(v)
		settingCache.Store(setting.CacheKey(), setting)
	}

}

func getEnvFile() (filename string, err error) {
	pwd, err := os.Getwd()
	if err != nil {
		return
	}
	envFileName := filepath.Join(pwd, ".env")
	var stat os.FileInfo
	if stat, err = os.Stat(envFileName); err == nil && !stat.IsDir() {
		filename = envFileName
		return
	}

	binPath, err := filepath.Abs(os.Args[0])
	if err != nil {
		return
	}
	envFileName = filepath.Join(filepath.Dir(binPath), ".env")
	if stat, err = os.Stat(envFileName); err == nil && !stat.IsDir() {
		filename = envFileName
		return
	}
	err = os.ErrNotExist
	return
}
