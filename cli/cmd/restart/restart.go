package restart

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"pmon3/pmond"
	"pmon3/pmond/model"
	"pmon3/pmond/output"
	"pmon3/pmond/svc/process"
)

var flag model.ExecFlags

var Cmd = &cobra.Command{
	Use:   "restart",
	Short: "restart some process by id or name",
	Run: func(cmd *cobra.Command, args []string) {
		cmdRun(args, flag.Json())
	},
}

func init() {
	Cmd.Flags().StringVarP(&flag.LogDir, "log_dir", "d", "", "the process stdout log dir")
	Cmd.Flags().StringVarP(&flag.Log, "log", "l", "", "the process stdout log")
}

func cmdRun(args []string, flags string) {
	if len(args) == 0 {
		pmond.Log.Fatal("please input restart process id or name")
	}

	val := args[0]
	var m model.Process
	if err := pmond.Db().First(&m, "id = ? or name = ?", val, val).Error; err != nil {
		pmond.Log.Fatal(fmt.Sprintf("the process %s not exist", val))
	}

	// checkout process state
	if process.IsRunning(m.Pid) {
		if err := process.TryStop(false, &m); err != nil {
			pmond.Log.Fatalf("restart error: %s", err.Error())
		}
	}

	rel, err := process.TryStart(m, flags)
	if err != nil {
		if len(os.Getenv("PMON3_DEBUG")) > 0 {
			pmond.Log.Fatalf("%+v", err)
		} else {
			pmond.Log.Fatal(err.Error())
		}
	}

	output.TableOne(rel)
}
