package conf

import (
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

// current app version
var Version = "1.14.3"

// two options:
// 1. Use PMON3_CONF envionment variable
// 2. fallback toa hardcoded path
func GetConfigFile() string {
	conf := os.Getenv("PMON3_CONF")
	if len(conf) == 0 {
		conf = "/etc/pmon3/config/config.yml"
	}
	return conf
}

func Load(configFile string) (*Tpl, error) {
	d, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var c Tpl
	err = yaml.Unmarshal(d, &c)
	if err != nil {
		return nil, err
	}

	c.ConfigFile = configFile

	return &c, nil
}
