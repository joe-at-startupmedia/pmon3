package conf

import "os"

// current app version
var Version = "1.13.0"

func GetDefaultConf() string {
	conf := os.Getenv("PMON3_CONF")
	if len(conf) == 0 {
		conf = "/etc/pmon3/config/config.yml"
	}
	return conf
}
