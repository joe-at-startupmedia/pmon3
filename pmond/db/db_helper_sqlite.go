//go:build !cqo_sqlite

package db

import (
	"github.com/joe-at-startupmedia/sqlite"
	"gorm.io/gorm"
)

func openDb(dbDir string) (*gorm.DB, error) {
	initDb, err := gorm.Open(sqlite.Open(dbDir+"/data.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	sqlDb, err := initDb.DB()
	sqlDb.SetMaxOpenConns(1)

	return initDb, err
}
