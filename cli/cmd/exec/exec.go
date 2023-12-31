package exec

import (
	"github.com/spf13/cobra"
	"os"
	"pmon2/pmond"
	"pmon2/pmond/model"
	"pmon2/pmond/output"
)

// process failed auto restart
var flag model.ExecFlags

var Cmd = &cobra.Command{
	Use:     "exec",
	Aliases: []string{"run"},
	Short:   "run one binary golang process file",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) <= 0 {
			pmond.Log.Fatal("porcess params is empty")
		}
		cmdRun(args, flag.Json())
	},
}

func init() {
	Cmd.Flags().BoolVarP(&flag.NoAutoRestart, "no-autorestart", "n", false, "not auto restart when process run failure")
	Cmd.Flags().StringVarP(&flag.User, "user", "u", "", "the process run user")
	Cmd.Flags().StringVarP(&flag.Log, "log", "l", "", "the process stdout log")
	Cmd.Flags().StringVarP(&flag.Args, "args", "a", "", "the process extra arguments")
	Cmd.Flags().StringVar(&flag.Name, "name", "", "run process name")
	Cmd.Flags().StringVarP(&flag.LogDir, "log_dir", "d", "", "the process stdout log dir")
}

func cmdRun(args []string, flags string) {
	// get exec abs file path
	execPath, err := getExecFile(args)
	if err != nil {
		pmond.Log.Error(err.Error())
		return
	}
	m, exist := processExist(execPath)
	var rel []string
	if exist {
		pmond.Log.Debugf("restart process: %v", flags)
		rel, err = restart(m, flags)
	} else {
		pmond.Log.Debugf("load first process: %v", flags)
		rel, err = loadFirst(execPath, flags)
	}

	if err != nil {
		if len(os.Getenv("PMON2_DEBUG")) > 0 {
			pmond.Log.Fatalf("%+v", err)
		} else {
			pmond.Log.Fatalf(err.Error())
		}
	}

	output.TableOne(rel)
}
