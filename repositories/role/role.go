package role

import (
	"github.com/kabukky/journey/conversion"
	"github.com/kabukky/journey/dao/scheme"
	"gorm.io/gorm"
)

func Create(db *gorm.DB, name string, description string, createByUser *scheme.User) (role *scheme.Role, err error) {
	role = &scheme.Role{}
	role.Name = conversion.XssFilter(name)
	role.Description = conversion.XssFilter(description)
	if createByUser != nil {
		role.CreatedBy = &createByUser.ID
	}

	err = db.Create(&role).Error
	if err != nil {
		role = nil
		return
	}
	return
}

func Update(db *gorm.DB, roleOrId interface{}, name string, description string, updateByUser *scheme.User) (role *scheme.Role, err error) {
	role, err = getRole(db, roleOrId)
	if err != nil {
		return
	}
	role.Name = conversion.XssFilter(name)
	role.Description = conversion.XssFilter(description)
	if updateByUser != nil {
		role.UpdatedBy = &updateByUser.ID
	}
	err = db.Save(role).Error
	return
}

func Delete(db *gorm.DB, roleOrId interface{}) (err error) {
	role, err := getRole(db, roleOrId)
	if err != nil {
		return
	}

	if err = RemoveAllUsers(db, role); err != nil {
		return
	}
	err = db.Delete(role).Error
	return
}

func List(db *gorm.DB) (roles []*scheme.Role, err error) {
	roles = make([]*scheme.Role, 0)
	err = db.Model(&scheme.Role{}).Find(&roles).Error
	if err != nil {
		return
	}
	return
}

func AddUsers(db *gorm.DB, role *scheme.Role, users ...*scheme.User) (err error) {
	err = db.Model(&role).Association("Users").Append(users)
	return
}

func RemoveAllUsers(db *gorm.DB, role *scheme.Role) (err error) {
	err = db.Model(&role).Association("Users").Clear()
	return
}

func ListUser(db *gorm.DB, role *scheme.Role) (users []*scheme.User, err error) {
	users = make([]*scheme.User, 0)
	err = db.Model(&role).Association("Users").Find(&users)
	if err != nil {
		users = nil
		return
	}
	return
}

func getRole(db *gorm.DB, roleOrId interface{}) (role *scheme.Role, err error) {
	if r, ok := roleOrId.(*scheme.Role); ok {
		role = r
	} else if searchName, ok := roleOrId.(string); ok {
		err = db.Model(&scheme.Role{}).Where("name = ?", searchName).First(&role).Error
		if err != nil {
			return
		}
	} else if searchId, ok := roleOrId.(uint); ok {
		err = db.Model(&scheme.Role{}).First(&role, searchId).Error
		if err != nil {
			return
		}
	}

	if role == nil {
		err = gorm.ErrRecordNotFound
	}

	return
}
