package tag

import (
	"github.com/kabukky/journey/dao/scheme"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

//func InsertTag(name []byte, slug string, created_at time.Time, created_by int64) (int64, error) {}
//func RetrieveTags(postId int64) ([]structure.Tag, error)                                        {}
//func RetrieveTag(tagId int64) (*structure.Tag, error)                                           {}
//func RetrieveTagBySlug(slug string) (*structure.Tag, error)                                     {}
//func RetrieveTagIdBySlug(slug string) (int64, error)                                            {}

var ErrTagIsExists = errors.New("tag is exists")
var ErrTagIsNotExists = errors.New("tag isn't exists")

func CreateTag(db *gorm.DB, tagName string, otherData map[string]interface{}) (ret *scheme.Tag, err error) {
	err = db.Model(&scheme.Tag{}).Where("name = ?", tagName).First(&ret).Error
	if err != nil && !errors.Is(gorm.ErrRecordNotFound, err) {
		return
	}

	if ret == nil {
		err = nil
	} else {
		return // tag is exists
	}

	ret = &scheme.Tag{}
	if otherData != nil {
		if err = mapstructure.Decode(otherData, ret); err != nil {
			return
		}
	}

	err = db.Create(&ret).Error
	return
}

func UpdateTag(db *gorm.DB, tagObjOrID interface{}, updateData map[string]interface{}) (tagObj *scheme.Tag, err error) {
	if t, ok := tagObjOrID.(*scheme.Tag); ok {
		tagObj = t
	} else if t, ok := tagObjOrID.(int64); ok {
		if err = db.Model(&scheme.Tag{}).First(&tagObj, t).Error; err != nil {
			return
		}
	}
	delete(updateData, "Id")
	delete(updateData, "UUID")
	if err = mapstructure.Decode(updateData, tagObj); err != nil {
		return
	}

	err = db.Save(tagObj).Error
	return
}

func DeleteTag(db *gorm.DB, tagObjOrID interface{}) (err error) {
	var tagObj *scheme.Tag

	if t, ok := tagObjOrID.(*scheme.Tag); ok {
		tagObj = t
	} else if t, ok := tagObjOrID.(int64); ok {
		if err = db.Model(&scheme.Tag{}).First(&tagObj, t).Error; err != nil {
			// always this tag is not exists
			if errors.Is(gorm.ErrRecordNotFound, err) {
				err = nil
			}
			return
		}
	}

	err = db.Delete(tagObj).Error
	return
}

func AllTag() (db *gorm.DB, tags scheme.Tags, err error) {
	tags = make(scheme.Tags, 0)
	err = db.Model(&scheme.Tag{}).Find(&tags).Error
	return
}

func GetTagByID(db *gorm.DB, tagID int64) (tag *scheme.Tag, err error) {
	err = db.Model(&scheme.Tag{}).Find(&tag, tagID).Error
	return
}

func GetTagBySlug(db *gorm.DB, slug string) (tag *scheme.Tag, err error) {
	err = db.Model(&scheme.Tag{}).Where("slug = ?", slug).First(&tag).Error
	return
}
