package templates

import (
	"github.com/kabukky/journey/dao/scheme"
	"github.com/kabukky/journey/helpers"
	"github.com/kabukky/journey/repositories/post"
	"github.com/kabukky/journey/repositories/tag"
	"github.com/kabukky/journey/repositories/user"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"net/http"
)

func ShowPostTemplate(c echo.Context, db *gorm.DB, slug string) (err error) {
	postObj, err := post.GetPostBySlug(db, slug)
	if err != nil {
		return
	}

	if postObj.Status != scheme.POST_STATUS_PUBLISHED {
		err = errors.New("Post not published.")
		return
	}

	tplCtx := map[string]interface{}{
		"Post":        postObj,
		"CurrentPath": c.Path(),
	}

	tplName := "post"
	if postObj.IsPage() {
		if _, err = GetTemplateFilePath("page"); err == nil {
			tplName = "page"
		}
		err = nil
	}

	htmlContent, err := Render(tplName, tplCtx)
	if err != nil {
		return
	}

	err = c.HTML(http.StatusOK, htmlContent)
	return
}

func ShowAuthorTemplate(c echo.Context, db *gorm.DB, slug string, page int) (err error) {
	author, err := user.GetUserBySlug(db, slug)
	if err != nil {
		return err
	}

	start, limit, err := helpers.Pagination(db, int64(page))

	conds := make(map[string]interface{})
	conds["user"] = author
	searchOpts := make(map[string]interface{})
	searchOpts["preload"] = true
	posts, total, err := post.GetPostBySearch(db, conds, start, limit, nil, searchOpts)

	if err != nil {
		return
	}

	tplCtx := map[string]interface{}{
		"Posts":       posts,
		"CurrentPath": c.Path(),
		"Total":       total,
		"Limit":       limit,
		"Start":       start,
	}

	tplName := "index"
	if _, err = GetTemplateFilePath("author"); err == nil {
		tplName = "author"
	}

	htmlContent, err := Render(tplName, tplCtx)
	if err != nil {
		return
	}

	err = c.HTML(http.StatusOK, htmlContent)
	return
}

func ShowTagTemplate(c echo.Context, db *gorm.DB, slug string, page int) (err error) {
	tagObj, err := tag.GetTagBySlug(db, slug)
	if err != nil {
		return err
	}

	start, limit, err := helpers.Pagination(db, int64(page))

	conds := make(map[string]interface{})
	conds["tags"] = tagObj.Name
	searchOpts := make(map[string]interface{})
	searchOpts["preload"] = true
	posts, total, err := post.GetPostBySearch(db, conds, start, limit, nil, searchOpts)

	if err != nil {
		return
	}

	tplCtx := map[string]interface{}{
		"Posts":       posts,
		"CurrentPath": c.Path(),
		"Total":       total,
		"Limit":       limit,
		"Start":       start,
	}

	tplName := "index"
	if _, err = GetTemplateFilePath("tag"); err == nil {
		tplName = "tag"
	}

	htmlContent, err := Render(tplName, tplCtx)
	if err != nil {
		return
	}

	err = c.HTML(http.StatusOK, htmlContent)
	return
}

func ShowIndexTemplate(c echo.Context, db *gorm.DB, page int) (err error) {
	start, limit, err := helpers.Pagination(db, int64(page))
	if err != nil {
		return
	}
	conds := make(map[string]interface{})
	if keyword := c.QueryParam("keyword"); keyword != "" {
		conds["keyword"] = keyword
	}
	searchOpts := make(map[string]interface{})
	searchOpts["preload"] = true
	posts, total, err := post.GetPostBySearch(db, conds, start, limit, nil, searchOpts)
	if err != nil {
		return
	}

	tplCtx := map[string]interface{}{
		"Posts":       posts,
		"CurrentPath": c.Path(),
		"Total":       total,
		"Limit":       limit,
		"Start":       start,
	}

	htmlContent, err := Render("index", tplCtx)
	if err != nil {
		return
	}

	err = c.HTML(http.StatusOK, htmlContent)
	return
}
