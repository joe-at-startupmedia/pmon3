package stop

import (
	"fmt"
	"os"
	"pmon3/pmond"
	"pmon3/pmond/model"
	"pmon3/pmond/output"
	process2 "pmon3/pmond/svc/process"

	"github.com/jinzhu/gorm"
	"github.com/spf13/cobra"
)

var (
	forceKill bool
)

var Cmd = &cobra.Command{
	Use:   "stop [id or name]",
	Short: "Stop a process by id or name",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cmdRun(args)
	},
}

func init() {
	Cmd.Flags().BoolVarP(&forceKill, "force", "f", false, "force the process to stop")
}

func cmdRun(args []string) {
	val := args[0]
	var process model.Process
	err := pmond.Db().Where("id = ? or name = ?", val, val).First(&process).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			pmond.Log.Fatalf("%s doesn't exist", val)
		}
	}

	StopProcess(&process, model.StatusStopped, forceKill)
}

func StopProcess(process *model.Process, status model.ProcessStatus, forced bool) error {
	// check process is running
	_, err := os.Stat(fmt.Sprintf("/proc/%d/status", process.Pid))
	if os.IsNotExist(err) {
		if process.Status == model.StatusRunning || process.Status != status {
			process.Status = status
			if err := pmond.Db().Save(&process).Error; err != nil {
				pmond.Log.Fatalf("stop process %s err \n", process.Stringify())
				return err
			}

			pmond.Log.Infof("stop process %s success \n", process.Stringify())
			return nil
		}
	}

	// try to kill the process
	err = process2.TryStop(process, status, forced)
	if err != nil {
		pmond.Log.Fatalf("stop the process %s failed", process.Stringify())
		return err
	}

	output.TableOne(process.RenderTable())
	return nil
}
