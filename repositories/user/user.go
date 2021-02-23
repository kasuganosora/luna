package user

import (
	"github.com/kabukky/journey/dao/scheme"
	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

/*
func InsertUser(name []byte, slug string, password string, email []byte, image []byte, cover []byte, created_at time.Time, created_by int64) (int64, error) {
}
func InsertRoleUser(role_id int, user_id int64) error                                            {}
func RetrieveUser(id int64) (*structure.User, error)                                             {}
func RetrieveUserBySlug(slug string) (*structure.User, error)                                    {}
func RetrieveUserByName(name []byte) (*structure.User, error)                                    {}
func RetrieveUsersCount() int                                                                    {}
func UpdateLastLogin(logInDate time.Time, userId int64) error                                    {}
func UpdateUserPassword(id int64, password string, updated_at time.Time, updated_by int64) error {}
*/

func CreateUser(db *gorm.DB, name string, password string, otherData map[string]interface{}) (user *scheme.User, err error) {
	user = &scheme.User{}
	if otherData != nil {
		delete(otherData, "UUID")
		delete(otherData, "ID")
		err = mapstructure.Decode(otherData, user)
		if err != nil {
			return
		}
	}

	if password != "" {
		if err = user.SetPassword(password); err != nil {
			return
		}
	}
	err = db.Create(&user).Error
	return
}

func UpdateUser(db gorm.DB, name string, updateData map[string]interface{}) (user *scheme.User, err error) {
	err = db.Model(&scheme.User{}).Where("name = ?", name).First(&user).Error
	if err != nil {
		return
	}
	if updateData != nil {
		delete(updateData, "UUID")
		delete(updateData, "ID")
		delete(updateData, "Password")
		delete(updateData, "PasswordSalt")
		delete(updateData, "CreatedAt")
		delete(updateData, "CreatedBy")
		err = mapstructure.Decode(updateData, user)
		if err != nil {
			return
		}
	}

	if passwd, ok := updateData["Password"]; ok && passwd.(string) != "" {
		err = user.SetPassword(passwd.(string))
		if err != nil {
			return
		}
	}

	err = db.Save(&user).Error
	return
}

func DeleteUser(db gorm.DB, userObjOrID interface{}) (err error) {
	var user *scheme.User
	if uid, ok := userObjOrID.(uint); ok {
		err = db.Model(&scheme.User{}).First(&user, uid).Error
		if err != nil {
			return
		}
	} else {
		user = userObjOrID.(*scheme.User)
	}

	err = db.Delete(user).Error
	return
}

func GetUserByID(db *gorm.DB, uid uint) (user scheme.User, err error) {
	err = db.Model(scheme.User{}).Preload(clause.Associations).First(&user, uid).Error
	return
}

func GetUserBySlug(db *gorm.DB, slug string) (user scheme.User, err error) {
	err = db.Model(scheme.User{}).Preload(clause.Associations).Where("slug = ?", slug).First(&user).Error
	return
}

func GetUserByName(db *gorm.DB, name string) (user scheme.User, err error) {
	err = db.Model(scheme.User{}).Preload(clause.Associations).Where("name = ?", name).First(&user).Error
	return
}

func UsersCount(db *gorm.DB) (count int64, err error) {
	err = db.Model(scheme.User{}).Count(&count).Error
	return
}
