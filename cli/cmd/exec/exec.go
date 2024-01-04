package exec

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"pmon3/cli/cmd/list"
	"pmon3/pmond"
	"pmon3/pmond/model"
	"time"

	"github.com/spf13/cobra"
)

// process failed auto restart
var flag model.ExecFlags

var Cmd = &cobra.Command{
	Use:     "exec",
	Aliases: []string{"run"},
	Short:   "Spawn a new process",
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
	parsedFlags, err := execflags.Parse(flags)
	if err != nil {
		pmond.Log.Fatalf("could not parse flags: %+v", err)
		return
	}
	name := parsedFlags.Name
	// get process file name
	if len(name) <= 0 {
		name = filepath.Base(args[0])
		parsedFlags.Name = name
	}
	err, _ = model.FindProcessByFileAndName(pmond.Db(), execPath, name)
	if err == nil {
		pmond.Log.Debugf("restart process: %v", flags)
		//@TODO moved to Restart controller
		//err = Restart(process, execPath, parsedFlags)
		err = loadFirst(execPath, parsedFlags)
		HandleStart(err)
	} else {
		pmond.Log.Debugf("load first process: %v", flags)
		err = loadFirst(execPath, parsedFlags)
		HandleStart(err)
	}
}

func HandleStart(err error) {
	if err != nil {
		if len(os.Getenv("PMON3_DEBUG")) > 0 {
			pmond.Log.Fatalf("%+v", err)
		} else {
			pmond.Log.Fatalf(err.Error())
		}
	}
	time.Sleep(pmond.Config.GetCmdExecResponseWait())
	list.Show()
}

// @TODO move to controller
func getExecFile(args []string) (string, error) {
	execFile := args[0]
	_, err := os.Stat(execFile)
	if os.IsNotExist(err) {
		return "", fmt.Errorf("%s does not exist", execFile)
	}

	if path.IsAbs(execFile) {
		return execFile, nil
	}

	absPath, err := filepath.Abs(execFile)
	if err != nil {
		return "", fmt.Errorf("get file path error: %v", err.Error())
	}

	return absPath, nil
}
