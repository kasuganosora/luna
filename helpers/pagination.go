package helpers

import (
	"github.com/kabukky/journey/repositories/setting"
	"gorm.io/gorm"
)

func Pagination(db *gorm.DB, currentPage int64, perPage ...int64) (start, limit int64, err error) {
	var ppage int64
	if len(perPage) > 0 {
		ppage = perPage[0]
	} else {
		ppage, err = setting.PostsPerPage(db, 15)
		if err != nil {
			return
		}
	}

	index := currentPage - 1
	if index < 0 {
		index = 0
	}
	start = index * ppage
	limit = ppage
	return
}
