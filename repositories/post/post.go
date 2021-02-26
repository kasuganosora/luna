package post

import (
	"github.com/elliotchance/pie/pie"
	"github.com/kabukky/journey/conversion"
	"github.com/kabukky/journey/dao"
	"github.com/kabukky/journey/dao/scheme"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
	"sync/atomic"
)

var ErrTagsTypeError = errors.New("tags type error")
var ErrPostNotExists = errors.New("post not exists")

var totalPostCount int64 = -1

func Create(db *gorm.DB, data map[string]interface{}) (post *scheme.Post, err error) {
	post = &scheme.Post{}
	if err = post.FillFromMap(data); err != nil {
		return
	}

	htmlData := conversion.GenerateHtmlFromMarkdown([]byte(post.Markdown))
	post.HTML = string(htmlData)

	if db == nil {
		db = dao.DB.Session(&gorm.Session{})
	}

	if err = db.Create(post).Error; err != nil {
		return
	}

	if tags, ok := data["tags_str"]; ok {
		if err = SetPostTags(db, post, tags, post.AuthorID); err != nil {
			return
		}
	}
	err = UpdateTotalPostCountCache(db)
	return
}

func Update(db *gorm.DB, postOrPostID interface{}, data map[string]interface{}) (post *scheme.Post, err error) {
	if db == nil {
		db = dao.DB.Session(&gorm.Session{})
	}

	post, err = getPostByIDOrSlug(db, postOrPostID)
	if err != nil {
		return
	}

	if err = post.FillFromMap(data); err != nil {
		return
	}

	htmlData := conversion.GenerateHtmlFromMarkdown([]byte(post.Markdown))
	post.HTML = string(htmlData)
	if err = db.Save(post).Error; err != nil {
		return
	}

	if tags, ok := data["tags_str"]; ok {
		if err = SetPostTags(db, post, tags, post.UpdatedBy); err != nil {
			return
		}
	}
	return
}

func Delete(db *gorm.DB, postOrPostID interface{}) (err error) {
	var post *scheme.Post
	if db == nil {
		db = dao.DB.Session(&gorm.Session{})
	}

	post, err = getPostByIDOrSlug(db, postOrPostID)
	if err != nil {
		return
	}

	err = db.Delete(post).Error
	if err == nil {
		err = UpdateTotalPostCountCache(db)
	}
	return
}

func SetPostTags(db *gorm.DB, post *scheme.Post, tags interface{}, setTagUserID *uint) (err error) {
	if tags == nil {
		return
	}

	if db == nil {
		db = dao.DB.Session(&gorm.Session{})
	}
	var _tags pie.Strings
	if _, ok := tags.([]string); ok {
		_tags = pie.Strings(tags.([]string))
	} else if v, ok := tags.(string); ok {
		_tags = strings.Split(v, ";")
	} else {
		err = ErrTagsTypeError
		return
	}

	tagsObj := make(scheme.Tags, 0)
	if err = db.Model(post).Association("Tags").Find(&tagsObj); err != nil {
		return
	}

	currentTags := tagsObj.StringsUsing(func(tag *scheme.Tag) string {
		return tag.Name
	})
	newTags, deleteTags := currentTags.Diff(_tags)
	optTags := make(scheme.Tags, 0)
	err = db.Model(&scheme.Tag{}).Where("name in ?", append(newTags, deleteTags...)).Find(&optTags).Error
	if err != nil {
		return
	}

	optTagsNameKeyMapping := make(map[string]*scheme.Tag)
	optTags.Each(func(tag *scheme.Tag) {
		optTagsNameKeyMapping[tag.Name] = tag
	})

	// add tags
	newTagsObj := make(scheme.Tags, len(newTags))
	for i, tag := range newTags {
		var tagObj *scheme.Tag
		var exists bool
		if tagObj, exists = optTagsNameKeyMapping[tag]; !exists {
			// tag not exists, now create this
			tagObj = &scheme.Tag{}
			tagObj.Name = tag
			if setTagUserID != nil {
				tagObj.CreatedBy = setTagUserID
			}
			if err = db.Create(tagObj).Error; err != nil {
				return
			}
			optTagsNameKeyMapping[tag] = tagObj
		}
		newTagsObj[i] = tagObj
	}

	if newTagsObj.Len() > 0 {
		if err = db.Model(post).Association("Tags").Append(newTagsObj); err != nil {
			return
		}
	}

	// remove tags
	removeTagsObj := make(scheme.Tags, 0)
	for _, tag := range deleteTags {
		if t, exists := optTagsNameKeyMapping[tag]; exists {
			removeTagsObj = append(removeTagsObj, t)
			continue
		}
	}

	if removeTagsObj.Len() > 0 {
		if err = db.Model(post).Association("Tags").Delete(removeTagsObj); err != nil {
			return
		}
	}
	return
}

func GetPostByID(db *gorm.DB, postID uint) (post *scheme.Post, err error) {
	tmp := scheme.Post{}
	err = db.Preload(clause.Associations).First(&tmp, postID).Error
	if err == nil {
		post = &tmp
	}
	return
}

func GetPostBySlug(db *gorm.DB, slug string) (post *scheme.Post, err error) {
	tmp := scheme.Post{}
	err = db.Preload(clause.Associations).Where("slug = ?", slug).First(&tmp).Error
	if err == nil {
		post = &tmp
	}
	return
}

func GetPostBySearch(db *gorm.DB, conditions map[string]interface{}, start, limit int64, orderBy interface{}, otherOpts map[string]interface{}) (posts scheme.Posts, total int64, err error) {
	query := db.Model(scheme.Post{})

	if conditions == nil {
		conditions = make(map[string]interface{})
	}

	if otherOpts == nil {
		otherOpts = make(map[string]interface{})
	}

	for key, val := range conditions {
		switch key {
		case "keyword":
			query = query.Where("(title like %?% or slug like %?%)", val.(string), val.(string))
		case "tags":
			var tags []string
			if v, ok := val.([]string); ok {
				tags = v
			}

			if v, ok := val.(string); ok {
				tags = strings.Split(v, ";")
			}

			if len(tags) > 0 {
				tagConds := db.Model(scheme.Tag{}).
					Joins("JOIN posts_tags ON posts_tags.tag_id = tags.id").
					Where("tags.name in ?", tags).
					Select("posts_tags.post_id")
				query = query.Where("id in (?)", tagConds)
			}

		case "slug":
			query = query.Where("slug = ?", val.(string))
		case "user":
			var targetUser *scheme.User
			if u, ok := val.(*scheme.User); ok {
				val = u
			}
			if s, ok := val.(string); ok {
				err = db.Model(&scheme.Post{}).Where("name = ?", s).First(&targetUser).Error
				if err != nil {
					if errors.Is(gorm.ErrRecordNotFound, err) {
						query = query.Where("1 != 1")
						break
					}
					return
				}
			}
			query = query.Where("author_id = ?", targetUser.ID)
		}
	}

	if err = query.Count(&total).Error; err != nil {
		return
	}

	if orderBy != nil {
		query = query.Order(orderBy)
	} else {
		query = query.Order("created_at DESC")
	}

	if limit > 0 {
		query = query.Limit(int(limit))
	}

	if start > 0 {
		query = query.Offset(int(start))
	}

	if b, ok := otherOpts["preload"]; ok && b.(bool) == true {
		query = query.Preload(clause.Associations)
	}
	err = query.Find(&posts).Error

	return
}

func UpdateTotalPostCountCache(db *gorm.DB) (err error) {
	var count int64
	err = db.Model(&scheme.Post{}).Count(&count).Error
	atomic.StoreInt64(&totalPostCount, count)
	return
}

func GetTotalPostCount(db *gorm.DB) (count int64, err error) {
	if totalPostCount == -1 {
		if err = UpdateTotalPostCountCache(db); err != nil {
			return
		}
	}

	count = totalPostCount
	return
}

func getPostByIDOrSlug(db *gorm.DB, idOrSlug interface{}) (post *scheme.Post, err error) {
	var ok bool
	tmp := scheme.Post{}
	if post, ok = idOrSlug.(*scheme.Post); ok {
		return
	}
	if id, ok := idOrSlug.(uint); ok {

		err = db.First(&tmp, id).Error
		if err != nil {
			post = nil
			return
		}
		post = &tmp
		return
	}

	if slug, ok := idOrSlug.(string); ok {
		err = db.Model(&scheme.Post{}).Where("slug = ?", slug).First(&tmp).Error
		if err != nil {
			post = nil
			return
		}
		post = &tmp
		return
	}
	post = nil
	err = ErrPostNotExists
	return
}
