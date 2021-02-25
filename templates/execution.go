package templates

import (
	"github.com/flosch/pongo2/v4"
	"github.com/kabukky/journey/dao"
	"github.com/kabukky/journey/filenames"
	"github.com/kabukky/journey/helpers"
	"github.com/kabukky/journey/repositories/setting"
	"os"
	"path/filepath"
)

func Render(template string, values map[string]interface{}) (content string, err error) {
	templateFilePath, err := GetTemplateFilePath(template)
	if err != nil {
		return
	}

	tpl, err := pongo2.FromFile(templateFilePath)
	if err != nil {
		return
	}

	if values == nil {
		values = pongo2.Context{}
	}

	blog, err := setting.RetrieveBlog(dao.DB)
	if err != nil {
		return
	}

	values["Blog"] = blog

	return tpl.Execute(values)
}

func GetTemplateFilePath(template string) (templateFilePath string, err error) {
	theme, err := setting.RetrieveActiveTheme(dao.DB)
	if err != nil {
		return
	}
	templateFilePath = filepath.Join(filenames.ThemesFilepath, theme.GetString(), template) + ".html"
	if _, err = os.Stat(templateFilePath); err != nil {
		templateFilePath = ""
		return
	}
	return
}

func GetAllThemes() []string {
	themes := make([]string, 0)
	files, _ := filepath.Glob(filepath.Join(filenames.ThemesFilepath, "*"))
	for _, file := range files {
		if helpers.IsDirectory(file) {
			themes = append(themes, filepath.Base(file))
		}
	}
	return themes
}
