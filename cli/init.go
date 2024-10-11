package cli

import (
	"pmon3/conf"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger
var Config *conf.Config

func Instance(confDir string) error {
	//cli doesnt need the process config file
	config, err := conf.Load(confDir, "")
	if err != nil {
		return err
	}
	Config = config
	Log = config.GetLogger()
	return nil
}
