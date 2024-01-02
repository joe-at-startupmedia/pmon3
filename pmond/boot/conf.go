package boot

import (
	"io/ioutil"
	"pmon3/pmond/conf"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

func Conf(confFile string) (*conf.Tpl, error) {
	d, err := ioutil.ReadFile(confFile)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var c conf.Tpl
	err = yaml.Unmarshal(d, &c)
	if err != nil {
		return nil, err
	}

	c.Conf = confFile

	return &c, nil
}
