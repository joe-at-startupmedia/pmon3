package drop

import (
	"pmon3/cli/cmd/del"
	"pmon3/pmond"
	"pmon3/pmond/model"

	"github.com/spf13/cobra"
)

var (
	forceKill bool
)

var Cmd = &cobra.Command{
	Use:   "drop",
	Short: "Delete all processes",
	Run: func(cmd *cobra.Command, args []string) {
		Drop()
	},
}

func init() {
	Cmd.Flags().BoolVarP(&forceKill, "force", "f", false, "force kill before deleting processes")
}

// show all process list
func Drop() {
	var all []model.Process
	err := pmond.Db().Find(&all).Error
	if err != nil {
		pmond.Log.Fatalf("pmon3 find process err: %v", err)
	}

	for _, process := range all {
		del.DelProcess(&process, forceKill)
	}
}
