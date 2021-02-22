package test

import (
	"github.com/kabukky/journey/dao/migration"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMigration(t *testing.T) {
	db, err := getDB()
	assert.Nil(t, err)

	err = migration.Do(db)
	assert.Nil(t, err)
}

func TestMigrationRollback(t *testing.T) {
	db, err := getDB()
	assert.Nil(t, err)

	err = migration.Do(db)
	assert.Nil(t, err)

	err = migration.Rollback(db)
	assert.Nil(t, err)
}
