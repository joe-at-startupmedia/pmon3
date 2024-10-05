//go:build cqo_sqlite

package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func openDb(dbFile string) (*gorm.DB, error) {
	return gorm.Open(sqlite.Open(dbFile), &gorm.Config{})
}
