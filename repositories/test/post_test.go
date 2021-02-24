package test

import (
	"github.com/kabukky/journey/dao/scheme"
	"github.com/kabukky/journey/repositories/post"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func TestPostCreate(t *testing.T) {
	db, err := getDB()
	assert.Nil(t, err)

	data := make(map[string]interface{})
	data["Title"] = "测试第一篇文章"
	data["Slug"] = "文章的Slug" + strconv.FormatInt(int64(rand.Intn(1000))+time.Now().Unix(), 10)
	data["Markdown"] = `
# H1  
## H2  
abc
`
	data["MetaTitle"] = "MetaTitle 标题"
	data["MetaDescription"] = "MetaDescription 内容"
	data["ScheduleTime"] = time.Now()
	data["tags_str"] = "test tag1;tag2;tag3"
	p, err := post.Create(db, data)
	assert.Nil(t, err)
	assert.NotNil(t, p)
}

func TestPostUpdate(t *testing.T) {
	db, err := getDB()
	assert.Nil(t, err)

	data := make(map[string]interface{})
	data["Title"] = "测试第一篇文章"
	data["Slug"] = "文章的Slug" + strconv.FormatInt(int64(rand.Intn(1000))+time.Now().Unix(), 10)
	data["Markdown"] = `
# H1  
## H2  
abc
`
	data["MetaTitle"] = "MetaTitle 标题"
	data["MetaDescription"] = "MetaDescription 内容"
	data["ScheduleTime"] = time.Now()
	data["tags_str"] = "test tag1;tag2;tag3"
	p, err := post.Create(db, data)
	assert.Nil(t, err)
	assert.NotNil(t, p)

	updateData := make(map[string]interface{})
	updateData["tags_str"] = "tag2;tag3;tag4"
	updateData["Title"] = "标题已修改"

	p, err = post.Update(db, p, updateData)
	assert.Nil(t, err)
	assert.NotNil(t, p)
}

func TestPostSearch(t *testing.T) {

	db, err := getDB()
	assert.Nil(t, err)

	// clear up test data
	deletePosts := make([]*scheme.Post, 0)
	db.Model(&scheme.Post{}).Where("title in ?", []string{"文章1", "文章2"}).Find(&deletePosts)
	if len(deletePosts) > 0 {
		for _, dp := range deletePosts {
			db.Delete(dp)
		}
	}

	data := make(map[string]interface{})
	data["Title"] = "文章1"
	data["Slug"] = "文章1Slug" + strconv.FormatInt(int64(rand.Intn(1000))+time.Now().Unix(), 10)
	data["Markdown"] = `
# H1  
## H2  
abc
`
	data["MetaTitle"] = "MetaTitle 标题"
	data["MetaDescription"] = "MetaDescription 内容"
	data["ScheduleTime"] = time.Now()
	data["tags_str"] = "test tag1;tag2;tag3"
	p, err := post.Create(db, data)
	assert.Nil(t, err)
	assert.NotNil(t, p)

	data["Title"] = "文章2"
	data["Slug"] = "文章2Slug" + strconv.FormatInt(int64(rand.Intn(1000))+time.Now().Unix(), 10)
	data["tags_str"] = "tag4"

	p, err = post.Create(db, data)
	assert.Nil(t, err)
	assert.NotNil(t, p)

	// test search
	// search all
	searchOpts := make(map[string]interface{})
	posts, total, err := post.GetPostBySearch(db, searchOpts, 0, 0, "id asc")
	assert.Nil(t, err)
	assert.True(t, total > 0)
	assert.NotNil(t, posts)
	assert.NotEmpty(t, posts)

	// search tag
	searchOpts = make(map[string]interface{})
	searchOpts["tags"] = "tag4"
	db = db.Debug()
	posts, total, err = post.GetPostBySearch(db, searchOpts, 0, 0, "id asc")
	assert.Nil(t, err)
	assert.True(t, total == 1)
	assert.NotNil(t, posts)
	assert.NotEmpty(t, posts)
}
