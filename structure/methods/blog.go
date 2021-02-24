package methods

import (
	"encoding/json"
	"github.com/kabukky/journey/configuration"
	"github.com/kabukky/journey/database"
	"github.com/kabukky/journey/date"
	"github.com/kabukky/journey/slug"
	"github.com/kabukky/journey/structure"
	"gorm.io/gorm"
	"log"
)

// Global blog - thread safe and accessible by all requests
var Blog *structure.Blog

var assetPath = []byte("/assets/")

func UpdateBlog(b *structure.Blog, userId int64) error {
	// Marshal navigation items to json string
	navigation, err := json.Marshal(b.NavigationItems)
	if err != nil {
		return err
	}
	err = database.UpdateSettings(b.Title, b.Description, b.Logo, b.Cover, b.PostsPerPage, b.ActiveTheme, navigation, date.GetCurrentTime(), userId)
	if err != nil {
		return err
	}
	// Generate new global blog
	err = GenerateBlog()
	if err != nil {
		log.Panic("Error: couldn't generate blog data:", err)
	}
	return nil
}

func UpdateActiveTheme(activeTheme string, userId int64) error {
	err := database.UpdateActiveTheme(activeTheme, date.GetCurrentTime(), userId)
	if err != nil {
		return err
	}
	// Generate new global blog
	err = GenerateBlog()
	if err != nil {
		log.Panic("Error: couldn't generate blog data:", err)
	}
	return nil
}

func GenerateBlog(db *gorm.DB) (err error) {

	return nil
}
