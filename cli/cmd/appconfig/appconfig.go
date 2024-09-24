package appconfig

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"pmon3/cli"
	"pmon3/cli/cmd/base"
)

var Cmd = &cobra.Command{
	Use:   "appconfig",
	Short: "Output current Application Configuration",
	Run: func(cmd *cobra.Command, args []string) {
		cmdRun(args)
	},
}

func jsonPrettyPrint(in string) string {
	var out bytes.Buffer
	err := json.Indent(&out, []byte(in), "", "  ")
	if err != nil {
		return in
	}
	return out.String()
}

func cmdRun(args []string) {
	base.OpenSender()
	defer base.CloseSender()
	sent := base.SendCmd("app_config", "")
	newCmdResp := base.GetResponse(sent)
	if len(newCmdResp.GetError()) > 0 {
		cli.Log.Fatalf(newCmdResp.GetError())
	}
	jsonOutput := newCmdResp.GetValueStr()
	fmt.Println(jsonPrettyPrint(jsonOutput))
}
