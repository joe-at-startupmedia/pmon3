package kill

import (
	"pmon3/pmond"
	"pmon3/pmond/model"

	"github.com/spf13/cobra"
)

var (
	forceKill bool
)

var Cmd = &cobra.Command{
	Use:   "kill",
	Short: "Terminate all processes",
	Run: func(cmd *cobra.Command, args []string) {
		Kill(model.StatusStopped)
	},
}

func init() {
	Cmd.Flags().BoolVarP(&forceKill, "force", "f", false, "force kill all processes")
}

// show all process list
func Kill(processStatus model.ProcessStatus) {
	var all []model.Process
	err := pmond.Db().Find(&all, "status = ?", model.StatusRunning).Error
	if err != nil {
		pmond.Log.Fatalf("pmon3 run err: %v", err)
	}
	/*
		for _, process := range all {
			stop.StopProcess(&process, processStatus, forceKill)
		}
	*/
}
