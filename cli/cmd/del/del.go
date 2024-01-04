package del

import (
	"fmt"
	"os"
	"pmon3/pmond"
	"pmon3/pmond/model"

	"github.com/spf13/cobra"
)

var (
	forceKill bool
)

var Cmd = &cobra.Command{
	Use:   "del [id or name]",
	Short: "Delete process by id or name",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		runCmd(args)
	},
}

func init() {
	Cmd.Flags().BoolVarP(&forceKill, "force", "f", false, "kill the process before deletion")
}

func runCmd(args []string) {
	val := args[0]
	var m model.Process
	err := pmond.Db().First(&m, "id = ? or name = ?", val, val).Error
	if err != nil {
		pmond.Log.Fatalf("del process err:%s \n", err.Error())
	}

	DelProcess(&m, forceKill)
	pmond.Log.Info("del process")
}

// show all process list
func DelProcess(process *model.Process, forceKill bool) {
	if process.Status == model.StatusRunning && !forceKill {
		pmond.Log.Fatalf(fmt.Sprintf("The process %s is running, you must must stop it first or pass the --force flag\n", process.Stringify()))
	}
	/*
		if forceKill {
			stop.StopProcess(process, model.StatusStopped, forceKill)
		}
	*/
	pmond.Db().Delete(process)
	_ = os.Remove(process.Log)
	pmond.Log.Info(fmt.Sprintf("Process %s successfully deleted", process.Stringify()))
}
