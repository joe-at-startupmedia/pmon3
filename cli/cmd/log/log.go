package log

import (
	"fmt"
	"github.com/spf13/cobra"
	"os/exec"
	"pmon3/pmond"
	"pmon3/pmond/model"
)

var Cmd = &cobra.Command{
	Use:   "log",
	Short: "display process log by id or name",
	Run: func(cmd *cobra.Command, args []string) {

		cmdRun(args)
	},
}

func cmdRun(args []string) {
	if len(args) == 0 {
		pmond.Log.Fatal("please input start process id or name")
	}
	val := args[0]
	var m model.Process
	if err := pmond.Db().First(&m, "id = ? or name = ?", val, val).Error; err != nil {
		pmond.Log.Fatal(fmt.Sprintf("the process %s not exist", val))
	}
	c := exec.Command("bash", "-c", "tail "+m.Log)
	output, _ := c.CombinedOutput()
	fmt.Println(string(output))

}
