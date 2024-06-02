package db

import (
	"os"
	"pmon3/conf"
	"pmon3/pmond"
	"pmon3/pmond/model"
	"sync"

	"gorm.io/gorm"
)

var dbOnce sync.Once
var db *gorm.DB

func Db() *gorm.DB {
	dbOnce.Do(func() {
		pmondDbDir := pmond.Config.DataDir
		_, err := os.Stat(pmondDbDir)
		if os.IsNotExist(err) {
			err := os.MkdirAll(pmondDbDir, 0755)
			if err != nil {
				pmond.Log.Panicf("%s", err)
			}
		}

		initDb, err := openDb(pmondDbDir)
		if err != nil {
			pmond.Log.Panicf("%s", err)
		}

		// init table
		if !initDb.Migrator().HasTable(&model.Process{}) {
			initDb.Migrator().CreateTable(&model.Process{})
		}

		if !initDb.Migrator().HasTable(&model.Pmond{}) {
			initDb.Migrator().CreateTable(&model.Pmond{})
		}

		// sync data
		var pmondModel model.Pmond
		err = initDb.First(&pmondModel).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound { // first version
				initDb.Create(&model.Pmond{Version: conf.Version})
			}
		}

		db = initDb
	})

	return db
}
