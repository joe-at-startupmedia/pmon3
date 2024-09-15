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

		db = initDb
	})

	return db
}
