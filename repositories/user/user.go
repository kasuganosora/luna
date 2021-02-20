package user

import (
	"github.com/kabukky/journey/structure"
	"time"
)

func InsertUser(name []byte, slug string, password string, email []byte, image []byte, cover []byte, created_at time.Time, created_by int64) (int64, error) {
}
func InsertRoleUser(role_id int, user_id int64) error                                            {}
func RetrieveUser(id int64) (*structure.User, error)                                             {}
func RetrieveUserBySlug(slug string) (*structure.User, error)                                    {}
func RetrieveUserByName(name []byte) (*structure.User, error)                                    {}
func RetrieveUsersCount() int                                                                    {}
func UpdateLastLogin(logInDate time.Time, userId int64) error                                    {}
func UpdateUserPassword(id int64, password string, updated_at time.Time, updated_by int64) error {}
