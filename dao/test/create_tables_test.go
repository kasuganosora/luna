package test

import (
	"github.com/kabukky/journey/dao/scheme"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateDB(t *testing.T) {
	db, err := getDB()
	assert.Nil(t, err)
	err = db.AutoMigrate(&scheme.Post{}, &scheme.Role{}, &scheme.User{}, &scheme.Tag{}, &scheme.Setting{})
	assert.Nil(t, err)
}
