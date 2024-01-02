package stop

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/spf13/cobra"
	"os"
	"pmon3/pmond"
	"pmon3/pmond/model"
	"pmon3/pmond/output"
	process2 "pmon3/pmond/svc/process"
)

var Cmd = &cobra.Command{
	Use:     "stop",
	Short:   "stop running process",
	Example: "sudo pmon3 stop [id or name]",
	Run: func(cmd *cobra.Command, args []string) {
		cmdRun(args)
	},
}

func cmdRun(args []string) {
	var val string
	if len(args) <= 0 {
		pmond.Log.Fatalf("must input some process id or name")
	}

	if len(args) == 1 {
		val = args[0]
	}

	// stop process force
	forced := false
	if len(args) == 2 {
		val = args[1]
		if args[0] == "-f" {
			forced = true
		}
	}

	var process model.Process
	err := pmond.Db().Where("id = ? or name = ?", val, val).First(&process).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			pmond.Log.Fatalf("%s not exist", val)
		}
	}

	// check process is running
	_, err = os.Stat(fmt.Sprintf("/proc/%d/status", process.Pid))
	if os.IsNotExist(err) {
		if process.Status == model.StatusRunning {
			process.Status = model.StatusStopped
			if err := pmond.Db().Save(&process).Error; err != nil {
				pmond.Log.Fatalf("stop process %s err \n", val)
			}

			pmond.Log.Infof("stop process %s success \n", val)
			return
		}
	}

	// try to kill the process
	err = process2.TryStop(forced, &process)
	if err != nil {
		pmond.Log.Fatalf("stop the process %s failed", val)
	}

	output.TableOne(process.RenderTable())
}
