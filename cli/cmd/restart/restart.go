package restart

import (
	"fmt"
	"pmon3/cli/cmd/exec"
	"pmon3/pmond"
	"pmon3/pmond/model"
	"pmon3/pmond/svc/process"

	"github.com/spf13/cobra"
)

var flag model.ExecFlags

var Cmd = &cobra.Command{
	Use:   "restart",
	Short: "Restart a process by id or name",
	Run: func(cmd *cobra.Command, args []string) {
		cmdRun(args, flag.Json())
	},
}

func init() {
	Cmd.Flags().BoolVarP(&flag.NoAutoRestart, "no-autorestart", "n", false, "do not restart upon process failure")
	Cmd.Flags().StringVarP(&flag.User, "user", "u", "", "the processes run user")
	Cmd.Flags().StringVarP(&flag.Log, "log", "l", "", "the processes stdout log")
	Cmd.Flags().StringVarP(&flag.Args, "args", "a", "", "the processes extra arguments")
	Cmd.Flags().StringVarP(&flag.LogDir, "log_dir", "d", "", "the processes stdout log dir")
}

func cmdRun(args []string, flags string) {
	if len(args) == 0 {
		pmond.Log.Fatal("please input restart process id or name")
	}

	idOrName := args[0]
	var m model.Process
	if err := pmond.Db().First(&m, "id = ? or name = ?", idOrName, idOrName).Error; err != nil {
		pmond.Log.Fatal(fmt.Sprintf("the process %s not exist", idOrName))
	}

	// checkout process state
	if process.IsRunning(m.Pid) {
		if err := process.TryStop(&m, model.StatusStopped, false); err != nil {
			pmond.Log.Fatalf("restart error: %s", err.Error())
		}
	}

	execflags := model.ExecFlags{}
	parsedFlags, err := execflags.Parse(flags)
	if err != nil {
		pmond.Log.Fatalf("could not parse flags: %+v", err)
		return
	}
	if err == nil {
		pmond.Log.Debugf("restart process: %v", flags)
		err = exec.Restart(&m, m.ProcessFile, parsedFlags)
		exec.HandleStart(err)
	}
}
