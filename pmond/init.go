package pmond

import (
	"os"
	"pmon3/conf"
	"pmon3/pmond/model"
	"sync"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger
var Config *conf.Tpl
var dbOnce sync.Once
var db *gorm.DB

func init() {
	Log = logrus.New()
	if os.Getenv("PMON3_DEBUG") == "true" {
		Log.SetLevel(logrus.DebugLevel)
		Log.SetReportCaller(true)
	} else {
		Log.SetLevel(logrus.InfoLevel)
	}
	Log.SetOutput(os.Stdout)
	Log.SetFormatter(&logrus.TextFormatter{
		DisableTimestamp: true,
	})
}

func Instance(confDir string) error {
	tpl, err := conf.Load(confDir)
	if err != nil {
		return err
	}

	Config = tpl

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

		initDb, err := gorm.Open("sqlite3", pmondDbDir+"/data.db")
		if err != nil {
			panic(err)
		}
		db = initDb

		// init table
		if !db.HasTable(&model.Process{}) {
			db.CreateTable(&model.Process{})
		}

		if !db.HasTable(&model.Pmond{}) {
			db.CreateTable(&model.Pmond{})
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
