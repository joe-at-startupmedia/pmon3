package db

import (
	"os"
	"pmon3/pmond"
	"sync"

	"gorm.io/gorm"
)

var dbOnce sync.Once
var db *gorm.DB

func Db() *gorm.DB {
	dbOnce.Do(func() {
		pmondDbDir := pmond.Config.Directory.Data
		_, err := os.Stat(pmondDbDir)
		if os.IsNotExist(err) {
			err = os.MkdirAll(pmondDbDir, 0755)
			if err != nil {
				pmond.Log.Panicf("%s", err)
			}
		}

		initDb, err := openDb(pmond.Config.GetDatabaseFile())
		if err != nil {
			pmond.Log.Panicf("%s", err)
		}

		db = initDb
	})

	return db
}
