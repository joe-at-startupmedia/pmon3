package list

import (
	"github.com/spf13/cobra"
	"pmon2/pmond"
	"pmon2/pmond/model"
	"pmon2/pmond/output"
)

var Cmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list"},
	Short:   "list all processes",
	Run: func(cmd *cobra.Command, args []string) {
		runCmd(nil)
	},
}

// show all process list
func runCmd(_ []string) {
	var all []model.Process
	err := pmond.Db().Find(&all).Error
	if err != nil {
		pmond.Log.Fatalf("pmon2 run err: %v", err)
	}

	var allProcess [][]string
	for _, process := range all {
		allProcess = append(allProcess, process.RenderTable())
	}

	output.Table(allProcess)
}
