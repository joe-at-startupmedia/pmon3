package proxy

import (
	"github.com/pkg/errors"
	"pmon2/cli/worker"
	"pmon2/pmond"
	conf2 "pmon2/pmond/conf"
	"pmon2/pmond/model"
	"pmon2/pmond/utils/array"
)

var cmdTypes = []string{"start", "restart"}

func RunProcess(args []string) ([]byte, error) {

	for _, arg := range args {
		pmond.Log.Debugf("RunProcess arg: %s\n", arg)
	}

	if len(args) <= 2 {
		return nil, errors.New("process params not valid")
	}
	conf := conf2.GetDefaultConf()
	err := pmond.Instance(conf)

	if err != nil {
		return nil, err
	}
	// check run type param
	typeCli := args[0]

	if !array.In(cmdTypes, typeCli) {
		return nil, errors.WithStack(err)
	}

	var output string

	flags := model.ExecFlags{}
	flagModel, err := flags.Parse(args[2])
	if err != nil {
		return nil, errors.WithStack(err)
	}

	switch typeCli {
	case "start":
		output, err = worker.Start(args[1], flagModel)
	case "restart":
		output, err = worker.Restart(args[1], flagModel)
	}

	if err != nil {
		return []byte(err.Error()), err
	}

	return []byte(output), nil
}