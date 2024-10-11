package db

import (
	"pmon3/pmond"
	"sync"

	"gorm.io/gorm"
)

var dbOnce sync.Once
var db *gorm.DB

func Db() *gorm.DB {
	dbOnce.Do(func() {
		pmondDbDir := pmond.Config.Data.Directory

		foc := pmond.Config.GetDataFileOwnershipConfig()

		err := foc.CreateDirectoryIfNonExistent(pmondDbDir)
		if err != nil {
			pmond.Log.Panicf("%s", err)
		}

		initDb, err := openDb(pmond.Config.GetDatabaseFile())
		if err != nil {
			pmond.Log.Panicf("%s", err)
		}

		foc.ApplyFilePermissions(pmond.Config.GetDatabaseFile())

		db = initDb
	})

	return db
}
