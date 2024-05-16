package conf

import (
	"encoding/json"
	"os"
	"pmon3/pmond/model"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

// current app version
var Version = "1.14.5"

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
	d, err := os.ReadFile(configFile)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var c Tpl
	err = yaml.Unmarshal(d, &c)
	if err != nil {
		return nil, err
	}

	c.ConfigFile = configFile

	if len(c.AppsConfigFile) > 0 {
		ac, err := LoadAppsJson(c.AppsConfigFile)
		if err != nil {
			return nil, err
		} else {
			c.AppsConfig = ac
		}
	}

	return &c, nil
}

type AppsConfig struct {
	Apps []AppsConfigApp `json:"apps"`
}

type AppsConfigApp struct {
	File  string          `json:"file"`
	Flags model.ExecFlags `json:"flags"`
}

func LoadAppsJson(configFile string) (*AppsConfig, error) {
	d, err := os.ReadFile(configFile)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var c AppsConfig
	err = json.Unmarshal(d, &c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}
