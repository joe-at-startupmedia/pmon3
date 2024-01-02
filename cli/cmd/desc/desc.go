package desc

import (
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/spf13/cobra"
	"pmon3/pmond"
	"pmon3/pmond/model"
	"pmon3/pmond/output"
	"strconv"
)

var Cmd = &cobra.Command{
	Use:     "desc",
	Aliases: []string{"show"},
	Short:   "print the process detail message",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			pmond.Log.Fatalf("The command need process name or id, Example: sudo pmon3 show test")
			return
		}

		cmdRun(args)
	},
}

func cmdRun(args []string) {
	val := args[0]

	var process model.Process
	err := pmond.Db().Find(&process, "name = ? or id = ?", val, val).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			pmond.Log.Fatalf("pmon3 run err: %v", err)
		}

		// not found
		pmond.Log.Errorf("process %s not exist", val)
		return
	}

	rel := [][]string{
		{"status", process.Status},
		{"id", strconv.Itoa(int(process.ID))},
		{"name", process.Name},
		{"pid", strconv.Itoa(process.Pid)},
		{"process", process.ProcessFile},
		{"args", process.Args},
		{"user", process.Username},
		{"log", process.Log},
		{"no-autorestart", process.NoAutoRestartStr()},
		{"created_at", process.CreatedAt.Format("2006-01-02 15:04:05")},
		{"updated_at", process.UpdatedAt.Format("2006-01-02 15:04:05")},
	}

	output.DescTable(rel)
}
