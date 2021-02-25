package server

import (
	"context"
	"fmt"
	"github.com/kabukky/journey/dao"
	"github.com/kabukky/journey/repositories/post"
	"github.com/kabukky/journey/repositories/setting"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"

	"github.com/kabukky/journey/filenames"
	"github.com/kabukky/journey/templates"
)

func indexHandler(c echo.Context) (err error) {
	number := c.Param("number")
	ctx := context.Background()
	db := dao.DB.WithContext(ctx)
	if number == "" {
		// Render index template (first page)
		err = templates.ShowIndexTemplate(c, db, 1)
		if err != nil {
			return
		}
		return
	}
	page, err := strconv.Atoi(number)
	if err != nil || page <= 1 {
		err = c.Redirect(http.StatusFound, "/")
		return
	}
	// Render index template
	err = templates.ShowIndexTemplate(c, db, page)
	if err != nil {
		return
	}
	return
}

func authorHandler(c echo.Context) (err error) {
	slug, _ := url.QueryUnescape(c.Param("slug"))
	function := c.Param("function")
	number := c.Param("number")

	ctx := context.Background()
	db := dao.DB.WithContext(ctx)

	if function == "" {
		// Render author template (first page)
		err = templates.ShowAuthorTemplate(c, db, slug, 1)
		if err != nil {
			return
		}
		return
	} else if function == "rss" {
		// Render author rss feed
		err = templates.ShowAuthorRss(c, db, slug)
		if err != nil {
			return
		}
		return
	}
	page, err := strconv.Atoi(number)
	if err != nil || page <= 1 {
		err = c.Redirect(http.StatusFound, "/")
		return
	}
	// Render author template
	err = templates.ShowAuthorTemplate(c, db, slug, page)
	if err != nil {
		return
	}
	return
}

func tagHandler(c echo.Context) (err error) {
	slug, _ := url.QueryUnescape(c.Param("slug"))
	function := c.Param("function")
	number := c.Param("number")

	ctx := context.Background()
	db := dao.DB.WithContext(ctx)

	if function == "" {
		// Render tag template (first page)
		err = templates.ShowTagTemplate(c, db, slug, 1)
		if err != nil {
			return
		}
		return
	} else if function == "rss" {
		// Render tag rss feed
		err = templates.ShowTagRss(c, db, slug)
		if err != nil {
			return
		}
		return
	}
	page, err := strconv.Atoi(number)
	if err != nil || page <= 1 {
		err = c.Redirect(http.StatusFound, "/")
		return
	}
	// Render tag template
	err = templates.ShowTagTemplate(c, db, slug, page)
	if err != nil {
		return
	}
	return
}

func postHandler(c echo.Context) (err error) {
	ctx := context.Background()
	db := dao.DB.WithContext(ctx)
	slug, _ := url.QueryUnescape(c.Param("slug"))
	if slug == "" {
		err = c.Redirect(http.StatusFound, "/")
		return
	} else if slug == "rss" {
		// Render index rss feed
		err = templates.ShowIndexRss(c, db)
		if err != nil {

			return
		}
		return
	}

	// Render post template
	err = templates.ShowPostTemplate(c, db, slug)
	if err != nil && err.Error() == "sql: no rows in result set" {
		http.Error(c.Response(), "Post Not found.", http.StatusNotFound)
		return
	}

	if err != nil {
		return
	}
	return
}

func postEditHandler(c echo.Context) (err error) {
	slug, _ := url.QueryUnescape(c.Param("slug"))

	if slug == "" {
		_ = c.Redirect(http.StatusFound, "/")
		return
	}

	ctx := context.Background()
	db := dao.DB.WithContext(ctx)

	postObj, err := post.GetPostBySlug(db, slug)
	if err != nil {
		c.Error(err)
		return
	}

	urlStr := fmt.Sprintf("/admin/#/edit/%d", postObj.ID)
	_ = c.Redirect(http.StatusTemporaryRedirect, urlStr)
	return
}

func faviconHandler(c echo.Context) (err error) {
	ex, err := os.Executable()
	if err != nil {
		return
	}
	iconPath := filepath.Join(filepath.Dir(ex), "favicon.ico")
	if _, err = os.Stat(iconPath); os.IsNotExist(err) {
		http.Error(c.Response(), "404 File not found.", http.StatusNotFound)
		return
	}

	file, err := os.OpenFile(iconPath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return
	}
	defer file.Close()

	_, _ = io.Copy(c.Response(), file)
	return
}

func InitializeBlog(router *echo.Echo) {
	router.GET("/", indexHandler)
	router.GET("/favicon.ico", faviconHandler)
	router.GET("/:slug/edit", postEditHandler)
	router.GET("/:slug/", postHandler)
	router.GET("/page/:number/", indexHandler)
	// For author
	router.GET("/author/:slug/", authorHandler)
	router.GET("/author/:slug/:function/", authorHandler)
	router.GET("/author/:slug/:function/:number/", authorHandler)
	// For tag
	router.GET("/tag/:slug/", tagHandler)
	router.GET("/tag/:slug/:function/", tagHandler)
	router.GET("/tag/:slug/:function/:number/", tagHandler)
	// For serving asset files
	themeSetting, err := setting.RetrieveActiveTheme(dao.DB)
	if err == nil {
		router.Static("/assets/", filepath.Join(filenames.ThemesFilepath, themeSetting.GetString(), "assets"))
	}
	router.Static("/images/", filenames.ImagesFilepath)
	router.Static("/content/images/", filenames.ImagesFilepath) // This is here to keep compatibility with Ghost
	router.Static("/public/", filenames.PublicFilepath)
}
