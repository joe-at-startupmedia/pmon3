//go:build cqo_sqlite

package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func openDb(dbDir string) (*gorm.DB, error) {
	return gorm.Open(sqlite.Open(dbDir+"/data.db"), &gorm.Config{})
}
