package test

import (
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
	data["Slug"] = "文章的Slug"
	data["Markdown"] = `
# H1  
## H2  
abc
`
	data["MetaTitle"] = "MetaTitle 标题"
	data["MetaDescription"] = "MetaDescription 内容"
	data["ScheduleTime"] = time.Now()
	data["tags_str"] = "test tag1;tag2;tag3"
	p, err := post.CreatePost(db, data)
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
	p, err := post.CreatePost(db, data)
	assert.Nil(t, err)
	assert.NotNil(t, p)

	updateData := make(map[string]interface{})
	updateData["tags_str"] = "tag2;tag3;tag4"
	updateData["Title"] = "标题已修改"

	p, err = post.UpdatePost(db, p, updateData)
	assert.Nil(t, err)
	assert.NotNil(t, p)
}
