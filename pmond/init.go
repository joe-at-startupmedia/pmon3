package pmond

import (
	"github.com/sirupsen/logrus"
	"pmon3/conf"
)

var Log *logrus.Logger
var Config *conf.Config

func Instance(confFile string, processConfFile string) error {
	Config = &conf.Config{}
	if err := conf.Load(confFile, processConfFile, Config); err != nil {
		return err
	}
	Log = Config.GetLogger()
	Log.Info(Config)
	return nil
}
