package cli

import (
	"pmon3/conf"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger
var Config *conf.Config

func Instance(confDir string) error {
	Config = &conf.Config{}
	//cli doesnt need the process config file
	if err := conf.Load(confDir, "", Config); err != nil {
		return err
	}
	Log = Config.GetLogger()
	return nil
}
