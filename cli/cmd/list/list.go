package list

import (
	"pmon3/pmond"
	"pmon3/pmond/model"
	"pmon3/pmond/output"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list"},
	Short:   "List all processes",
	Run: func(cmd *cobra.Command, args []string) {
		Show()
	},
}

// show all process list
func Show() {
	var all []model.Process
	err := pmond.Db().Find(&all).Error
	if err != nil {
		pmond.Log.Fatalf("pmon3 run err: %v", err)
	}

	var allProcess [][]string
	for _, process := range all {
		allProcess = append(allProcess, process.RenderTable())
	}

	output.Table(allProcess)
}
