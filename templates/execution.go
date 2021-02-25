package templates

import (
	"bytes"
	"errors"
	"github.com/kabukky/journey/dao/scheme"
	"github.com/kabukky/journey/filenames"
	"github.com/kabukky/journey/helpers"
	"github.com/kabukky/journey/repositories/post"
	"github.com/kabukky/journey/repositories/setting"
	tag2 "github.com/kabukky/journey/repositories/tag"
	"github.com/kabukky/journey/repositories/user"
	"github.com/kabukky/journey/structure"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"path/filepath"
	"sync"
)

type Templates struct {
	sync.RWMutex
	m map[string]*structure.Helper
}

func newTemplates() *Templates { return &Templates{m: make(map[string]*structure.Helper)} }

// Global compiled templates - thread safe and accessible by all requests
var compiledTemplates = newTemplates()

func ShowPostTemplate(c echo.Context, db *gorm.DB, slug string) (err error) {
	// Read lock templates and global blog
	compiledTemplates.RLock()
	defer compiledTemplates.RUnlock()

	postObj, err := post.GetPostBySlug(db, slug)
	if err != nil {
		return err
	} else if postObj.PublishedAt == nil { // Make sure the post is published before rendering it
		return errors.New("Post not published.")
	}

	blog, err := setting.RetrieveBlog(db)
	if err != nil {
		return
	}

	requestData := structure.RequestData{Posts: make([]*scheme.Post, 1), Blog: blog, CurrentTemplate: 1, CurrentPath: c.Path()} // CurrentTemplate = post
	requestData.Posts[0] = postObj
	// Check if there's a custom page template available for this slug
	if template, ok := compiledTemplates.m["page-"+slug]; ok {
		_, err = c.Response().Write(executeHelper(template, &requestData, 1)) // context = post
		return err
	}
	// If the post is a page and the page template is available, use the page template
	if postObj.Page {
		if template, ok := compiledTemplates.m["page"]; ok {
			_, err = c.Response().Write(executeHelper(template, &requestData, 1)) // context = post
			return err
		}
	}
	_, err = c.Response().Write(executeHelper(compiledTemplates.m["post"], &requestData, 1)) // context = post

	return err
}

func ShowAuthorTemplate(c echo.Context, db *gorm.DB, slug string, page int) (err error) {
	// Read lock templates and global blog
	compiledTemplates.RLock()
	defer compiledTemplates.RUnlock()

	postIndex := int64(page - 1)
	if postIndex < 0 {
		postIndex = 0
	}

	postsPerPage, err := setting.PostsPerPage(db, 10)
	if err != nil {
		return
	}

	author, err := user.GetUserBySlug(db, slug)
	if err != nil {
		return err
	}

	conds := make(map[string]interface{})
	conds["user"] = author
	posts, _, err := post.GetPostBySearch(db, conds, (postsPerPage * postIndex), postsPerPage, nil)

	if err != nil {
		return err
	}

	blog, err := setting.RetrieveBlog(db)
	if err != nil {
		return
	}

	requestData := structure.RequestData{Posts: posts, Blog: blog, CurrentIndexPage: page, CurrentTemplate: 3, CurrentPath: c.Path()} // CurrentTemplate = author
	if template, ok := compiledTemplates.m["author"]; ok {
		_, err = c.Response().Write(executeHelper(template, &requestData, 0)) // context = index
	} else {
		_, err = c.Response().Write(executeHelper(compiledTemplates.m["index"], &requestData, 0)) // context = index
	}

	return err
}

func ShowTagTemplate(c echo.Context, db *gorm.DB, slug string, page int) (err error) {
	// Read lock templates and global blog
	compiledTemplates.RLock()
	defer compiledTemplates.RUnlock()

	postIndex := int64(page - 1)
	if postIndex < 0 {
		postIndex = 0
	}
	tag, err := tag2.GetTagBySlug(db, slug)
	if err != nil {
		return err
	}

	postsPerPage, err := setting.PostsPerPage(db, 10)
	if err != nil {
		return
	}

	conds := make(map[string]interface{})
	conds["tags"] = tag.Name
	posts, _, err := post.GetPostBySearch(db, conds, (postsPerPage * postIndex), postsPerPage, nil)

	if err != nil {
		return err
	}

	blog, err := setting.RetrieveBlog(db)
	if err != nil {
		return
	}

	requestData := structure.RequestData{Posts: posts, Blog: blog, CurrentIndexPage: page, CurrentTag: tag, CurrentTemplate: 2, CurrentPath: c.Path()} // CurrentTemplate = tag
	if template, ok := compiledTemplates.m["tag"]; ok {
		_, err = c.Response().Write(executeHelper(template, &requestData, 0)) // context = index
	} else {
		_, err = c.Response().Write(executeHelper(compiledTemplates.m["index"], &requestData, 0)) // context = index
	}

	return err
}

func ShowIndexTemplate(c echo.Context, db *gorm.DB, page int) (err error) {
	// Read lock templates and global blog
	compiledTemplates.RLock()
	defer compiledTemplates.RUnlock()

	postIndex := int64(page - 1)
	if postIndex < 0 {
		postIndex = 0
	}

	conds := make(map[string]interface{})

	postsPerPage, err := setting.PostsPerPage(db, 10)
	if err != nil {
		return
	}

	posts, _, err := post.GetPostBySearch(db, conds, (postsPerPage * postIndex), postsPerPage, nil)
	if err != nil {
		return err
	}

	blog, err := setting.RetrieveBlog(db)
	if err != nil {
		return
	}

	requestData := structure.RequestData{Posts: posts, Blog: blog, CurrentIndexPage: page, CurrentTemplate: 0, CurrentPath: c.Path()} // CurrentTemplate = index

	_, err = c.Response().Write(executeHelper(compiledTemplates.m["index"], &requestData, 0)) // context = index
	return err
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

func executeHelper(helper *structure.Helper, values *structure.RequestData, context int) []byte {
	// Set context and set it back to the old value once fuction returns
	defer setCurrentHelperContext(values, values.CurrentHelperContext)
	values.CurrentHelperContext = context

	block := helper.Block
	indexTracker := 0
	extended := false
	var extendHelper *structure.Helper
	for index, child := range helper.Children {
		// Handle extend helper
		if index == 0 && child.Name == "!<" {
			extended = true
			extendHelper = compiledTemplates.m[string(child.Function(&child, values))]
		} else {
			var buffer bytes.Buffer
			toAdd := child.Function(&child, values)
			buffer.Write(block[:child.Position+indexTracker])
			buffer.Write(toAdd)
			buffer.Write(block[child.Position+indexTracker:])
			block = buffer.Bytes()
			indexTracker += len(toAdd)
		}
	}
	if extended {
		extendHelper.BodyHelper.Block = block
		return executeHelper(extendHelper, values, values.CurrentHelperContext) // TODO: not sure if context = values.CurrentHelperContext is right.
	}
	return block
}

func setCurrentHelperContext(values *structure.RequestData, context int) {
	values.CurrentHelperContext = context
}
