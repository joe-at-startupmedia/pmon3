package exec

import (
	"os"
	"path/filepath"
	"pmon3/pmond"
	"pmon3/pmond/model"
	"pmon3/pmond/output"

	"github.com/spf13/cobra"
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
	Cmd.Flags().BoolVarP(&flag.NoAutoRestart, "no-autorestart", "n", false, "do not restart upon process failure")
	Cmd.Flags().StringVarP(&flag.User, "user", "u", "", "the processes run user")
	Cmd.Flags().StringVarP(&flag.Log, "log", "l", "", "the processes stdout log")
	Cmd.Flags().StringVarP(&flag.Args, "args", "a", "", "the processes extra arguments")
	Cmd.Flags().StringVar(&flag.Name, "name", "", "the processes name")
	Cmd.Flags().StringVarP(&flag.LogDir, "log_dir", "d", "", "the processes stdout log dir")
}

func cmdRun(args []string, flags string) {
	// get exec abs file path
	execPath, err := getExecFile(args)
	if err != nil {
		pmond.Log.Error(err.Error())
		return
	}
	execflags := model.ExecFlags{}
	flagModel, err := execflags.Parse(flags)
	if err != nil {
		pmond.Log.Fatalf("could not parse flags: %+v", err)
		return
	}
	name := flagModel.Name
	// get process file name
	if len(name) <= 0 {
		name = filepath.Base(args[0])
	}
	err, process := model.FindByProcessFileAndName(pmond.Db(), execPath, name)
	var rel []string
	if err == nil {
		pmond.Log.Debugf("restart process: %v", flags)
		rel, err = restart(process, flags)
	} else {
		pmond.Log.Debugf("load first process: %v", flags)
		rel, err = loadFirst(execPath, flags)
	}

	if err != nil {
		if len(os.Getenv("PMON3_DEBUG")) > 0 {
			pmond.Log.Fatalf("%+v", err)
		} else {
			pmond.Log.Fatalf(err.Error())
		}
	}

	output.TableOne(rel)
}
