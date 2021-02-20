package scripts

import (
	"github.com/kabukky/journey/dao/scheme"
	"gorm.io/gorm"
)

type InstallScript struct {
}

func (s *InstallScript) Do(db *gorm.DB) (err error) {
	// create Tables
	err = db.AutoMigrate(&scheme.Post{}, &scheme.Role{}, &scheme.User{}, &scheme.Tag{}, &scheme.Setting{})
	if err != nil {
		return
	}
	settingRecords := make([]scheme.Setting, 0)
	settingRecords = append(settingRecords, scheme.Setting{Key: "Title", Value: "My Blog", Type: "blog"})
	settingRecords = append(settingRecords, scheme.Setting{Key: "description", Value: "Just another Blog", Type: "blog"})
	settingRecords = append(settingRecords, scheme.Setting{Key: "logo", Value: "/public/images/blog-logo.jpg", Type: "blog"})
	settingRecords = append(settingRecords, scheme.Setting{Key: "cover", Value: "/public/images/blog-cover.jpg", Type: "blog"})
	settingRecords = append(settingRecords, scheme.Setting{Key: "postsPerPage", Value: "5", Type: "blog"})
	settingRecords = append(settingRecords, scheme.Setting{Key: "activeTheme", Value: "promenade", Type: "blog"})
	settingRecords = append(settingRecords, scheme.Setting{Key: "navigation", Value: `[{"label":"Home", "url":"/"}]`, Type: "blog"})

	for _, record := range settingRecords {
		r := &scheme.Setting{}
		err = db.Where("key = ?", record.Key).First(&r).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return // unknown error
		}
		if err == nil {
			continue // setting is exists
		}
		// new setting

		err = db.Create(&record).Error
		if err != nil {
			return
		}
	}
	return
}
func (s *InstallScript) Rollback(db *gorm.DB) (err error) {
	err = db.Migrator().DropTable(
		&scheme.Post{},
		&scheme.Role{},
		&scheme.User{},
		&scheme.Tag{},
		&scheme.Setting{},
	)
	return
}
func (s *InstallScript) Name() string {
	return "install_script"
}
