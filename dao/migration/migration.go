package migration

import (
	"github.com/elliotchance/pie/pie"
	"github.com/kabukky/journey/dao/scheme"
	"gorm.io/gorm"
)

var scriptsName []string
var scriptNameMapping map[string]ScriptInterface

type ScriptInterface interface {
	Do(db *gorm.DB) (err error)
	Rollback(db *gorm.DB) (err error)
	Name() string
}

func init() {
	scriptsName = make([]string, len(scripts))
	scriptNameMapping = make(map[string]ScriptInterface)
	for i, s := range scripts {
		scriptsName[i] = s.Name()
		scriptNameMapping[s.Name()] = s
	}
}

func initMigration(db *gorm.DB) error {
	return db.AutoMigrate(&scheme.Migration{})
}

func Do(db *gorm.DB) (err error) {
	if err = initMigration(db); err != nil {
		return
	}

	var oldScripts []string
	err = db.Model(&scheme.Migration{}).Pluck("name", &oldScripts).Error
	if err != nil {
		return
	}

	newScripts, _ := pie.Strings(oldScripts).Diff(scriptsName)
	if newScripts.Len() == 0 {
		return
	}

	// get last Batch num
	lastBatch, err := getLastBatchNum(db)
	if err != nil {
		return
	}

	lastBatch += 1

	for _, n := range newScripts {
		script := scriptNameMapping[n]
		err = script.Do(db)
		if err != nil {
			return
		}
		// add migration record
		record := scheme.Migration{
			Name:  n,
			Batch: lastBatch,
		}
		err = db.Create(&record).Error
		if err != nil {
			return
		}
	}
	return
}

func Rollback(db *gorm.DB) (err error) {
	if err = initMigration(db); err != nil {
		return
	}
	lastBatchNum, err := getLastBatchNum(db)
	if err != nil || lastBatchNum == 0 {
		return
	}

	// get last batch scripts
	var scriptsName []string
	err = db.Model(&scheme.Migration{}).
		Where("batch = ?", lastBatchNum).
		Order("id DESC").
		Pluck("name", &scriptsName).
		Error
	if err != nil {
		return
	}

	for _, n := range scriptsName {
		script := scriptNameMapping[n]
		err = script.Do(db)
		if err != nil {
			return
		}
		err = script.Rollback(db)
		if err != nil {
			return
		}
		err = db.
			Where("name = ?", n).
			Where("batch = ?", lastBatchNum).
			Delete(&scheme.Migration{}).
			Error
		if err != nil {
			return
		}
	}
	return
}

func getLastBatchNum(db *gorm.DB) (lastBatch int64, err error) {
	lastRecord := &scheme.Migration{}

	if err = db.Model(&scheme.Migration{}).Order("batch desc").First(lastRecord).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			lastBatch = 0
			err = nil
			return
		} else {
			return
		}
	} else {
		lastBatch = lastRecord.Batch
	}

	return
}
