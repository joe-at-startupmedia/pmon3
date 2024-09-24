package export

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"pmon3/cli"
	"pmon3/cli/cmd/base"
	"pmon3/pmond/model"
)

type flags struct {
	format string
}

var flag flags

var Cmd = &cobra.Command{
	Use:   "export",
	Short: "Export Process Configuration",
	Run: func(cmd *cobra.Command, args []string) {
		cmdRun(args)
	},
}

func init() {
	Cmd.Flags().StringVarP(&flag.format, "format", "f", "json", "the format to export")
}

func jsonPrettyPrint(in string) string {
	var out bytes.Buffer
	err := json.Indent(&out, []byte(in), "", "  ")
	if err != nil {
		return in
	}
	return out.String()
}

func tomlPrettyPrint(ac *model.ProcessConfig) string {
	out := new(bytes.Buffer)
	encoder := toml.NewEncoder(out)
	if err := encoder.Encode(ac); err != nil {
		cli.Log.Fatal(err)
	}
	return out.String()
}

func yamlPrettyPrint(ac *model.ProcessConfig) string {
	out := new(bytes.Buffer)
	encoder := yaml.NewEncoder(out)
	if err := encoder.Encode(ac); err != nil {
		cli.Log.Fatal(err)
	}
	return out.String()
}

func cmdRun(args []string) {
	base.OpenSender()
	defer base.CloseSender()
	sent := base.SendCmd("export", "")
	newCmdResp := base.GetResponse(sent)
	if len(newCmdResp.GetError()) > 0 {
		cli.Log.Fatalf(newCmdResp.GetError())
	}
	jsonOutput := newCmdResp.GetValueStr()

	var ac model.ProcessConfig
	if flag.format == "toml" || flag.format == "yaml" {
		err := json.Unmarshal([]byte(jsonOutput), &ac)
		if err != nil {
			cli.Log.Fatal(err)
		}
	}

	switch flag.format {
	case "toml":
		fmt.Println(tomlPrettyPrint(&ac))
	case "yaml":
		fmt.Println(yamlPrettyPrint(&ac))
	case "json":
		fmt.Println(jsonPrettyPrint(jsonOutput))
	default:
		cli.Log.Fatalf("Formats accepted are: json, toml or yaml")
	}
}
