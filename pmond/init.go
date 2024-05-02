package pmond

import (
	"os"
	"pmon3/conf"
	"pmon3/pmond/model"
	"sync"

	"github.com/joe-at-startupmedia/sqlite"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var Log *logrus.Logger
var Config *conf.Tpl
var dbOnce sync.Once
var db *gorm.DB

func Instance(confDir string) error {
	tpl, err := conf.Load(confDir)
	if err != nil {
		return err
	}

	Config = tpl

	Log = logrus.New()
	loglevel := tpl.GetLogrusLevel()
	if loglevel > logrus.WarnLevel {
		Log.SetReportCaller(true)
	}
	Log.SetLevel(loglevel)
	Log.SetOutput(os.Stdout)
	Log.SetFormatter(&logrus.TextFormatter{
		DisableTimestamp: true,
	})

	return nil
}

func Db() *gorm.DB {
	dbOnce.Do(func() {
		pmondDbDir := Config.Data
		_, err := os.Stat(pmondDbDir)
		if os.IsNotExist(err) {
			err := os.MkdirAll(pmondDbDir, 0755)
			if err != nil {
				panic(err)
			}
		}

		initDb, err := gorm.Open(sqlite.Open(pmondDbDir+"/data.db"), &gorm.Config{})
		if err != nil {
			panic(err)
		}
		db = initDb

		// init table
		if !db.Migrator().HasTable(&model.Process{}) {
			db.Migrator().CreateTable(&model.Process{})
		}

		if !db.Migrator().HasTable(&model.Pmond{}) {
			db.Migrator().CreateTable(&model.Pmond{})
		}

		// sync data
		var pmondModel model.Pmond
		err = db.First(&pmondModel).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound { // first version
				db.Create(&model.Pmond{Version: conf.Version})
			}
		}
	})

	return db
}
