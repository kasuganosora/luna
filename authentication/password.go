package authentication

import (
	"github.com/kabukky/journey/dao"
	"github.com/kabukky/journey/logger"
	"github.com/kabukky/journey/repositories/user"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func LoginIsCorrect(name string, password string) bool {
	userObj, err := user.GetUserByName(dao.DB, name)
	if err != nil && !errors.Is(gorm.ErrRecordNotFound, err) {
		logger.Error("LoginIsCorrect has error: %v", err)
	}

	ok, err := userObj.ComparePassword(password)
	if err != nil {
		logger.Error("LoginIsCorrect has error: %v", err)
	}

	return ok
}
