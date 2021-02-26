package templates

import (
	"bytes"
	"github.com/kabukky/feeds"
	"github.com/kabukky/journey/date"
	"github.com/kabukky/journey/repositories/post"
	"github.com/kabukky/journey/repositories/setting"
	tag2 "github.com/kabukky/journey/repositories/tag"
	"github.com/kabukky/journey/repositories/user"
	"github.com/kabukky/journey/structure"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func ShowIndexRss(c echo.Context, db *gorm.DB) (err error) {
	// Read lock global blog

	// 15 posts in rss for now
	//posts, err := database.RetrievePostsForIndex(15, 0)

	blog, err := setting.RetrieveBlog(db)
	if err != nil {
		return
	}

	searchOpts := make(map[string]interface{})
	searchOpts["preload"] = true

	posts, _, err := post.GetPostBySearch(db, nil, 0, 15, "created_at desc", searchOpts)
	if err != nil {
		return err
	}
	blogData := &structure.RequestData{Posts: posts, Blog: blog}
	feed := createFeed(blogData)
	err = feed.WriteRss(c.Response())
	return err
}

func ShowTagRss(c echo.Context, db *gorm.DB, slug string) (err error) {
	// Read lock global blog

	blog, err := setting.RetrieveBlog(db)
	if err != nil {
		return
	}

	tagObj, err := tag2.GetTagBySlug(db, slug)
	if err != nil {
		return
	}
	// 15 posts in rss for now
	conds := make(map[string]interface{})
	conds["tags"] = tagObj.Name
	searchOpts := make(map[string]interface{})
	searchOpts["preload"] = true
	posts, _, err := post.GetPostBySearch(db, conds, 0, 15, "created_at desc", searchOpts)
	if err != nil {
		return
	}
	blogData := &structure.RequestData{Posts: posts, Blog: blog}
	feed := createFeed(blogData)
	err = feed.WriteRss(c.Response())
	return err
}

func ShowAuthorRss(c echo.Context, db *gorm.DB, slug string) (err error) {
	blog, err := setting.RetrieveBlog(db)
	if err != nil {
		return
	}

	author, err := user.GetUserBySlug(db, slug)
	if err != nil {
		return
	}

	// 15 posts in rss for now
	conds := make(map[string]interface{})
	conds["user"] = author.ID
	searchOpts := make(map[string]interface{})
	searchOpts["preload"] = true
	posts, _, err := post.GetPostBySearch(db, conds, 0, 15, "created_at desc", searchOpts)
	if err != nil {
		return
	}
	blogData := &structure.RequestData{Posts: posts, Blog: blog}
	feed := createFeed(blogData)
	err = feed.WriteRss(c.Response())
	return err
}

func createFeed(values *structure.RequestData) *feeds.Feed {
	now := date.GetCurrentTime()
	feed := &feeds.Feed{
		Title:       string(values.Blog.Title),
		Description: string(values.Blog.Description),
		Link:        &feeds.Link{Href: string(values.Blog.Url)},
		Updated:     now,
		Image: &feeds.Image{
			Url:   string(values.Blog.Url) + string(values.Blog.Logo),
			Title: string(values.Blog.Title),
			Link:  string(values.Blog.Url),
		},
		Url: string(values.Blog.Url) + "/rss/",
	}

	for _, postObj := range values.Posts {
		if postObj.ID <= 0 {
			continue
		}
		// Make link
		var buffer bytes.Buffer
		buffer.WriteString(values.Blog.Url)
		buffer.WriteString("/")
		if postObj.Slug == nil {
			continue
		}
		buffer.WriteString(*postObj.Slug)
		item := &feeds.Item{
			Title:       postObj.Title,
			Description: postObj.HTML,
			Link:        &feeds.Link{Href: buffer.String()},
			Id:          postObj.UUID.String(),
			Author:      &feeds.Author{Name: postObj.Author.Name, Email: ""},
			Created:     postObj.CreatedAt,
		}

		image := string(postObj.Image)
		if image != "" {
			item.Image = &feeds.Image{
				Url: string(values.Blog.Url) + image,
			}
		}
		feed.Items = append(feed.Items, item)

	}

	return feed
}
